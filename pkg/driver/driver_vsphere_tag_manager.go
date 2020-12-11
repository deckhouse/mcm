package driver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

type TagManager struct {
	sync.Mutex

	regionTagCategoryName string
	zoneTagCategoryName   string

	ensuredTagCategoriesIDs      map[string]string
	ensuredTagIDsForCategoryName map[string]map[string]string

	regionTagNameToTag map[string]*tags.Tag
	zoneTagNameToTag   map[string]*tags.Tag

	tagIdToDcRef      map[string]object.Reference
	tagIdToClusterRef map[string]object.Reference
}

func NewTagManager(regionTagCategoryName, zoneTagCategoryName string) *TagManager {
	tm := &TagManager{
		regionTagCategoryName: regionTagCategoryName,
		zoneTagCategoryName:   zoneTagCategoryName,
	}

	go wait.Forever(func() {
		tm.Lock()
		defer tm.Unlock()

		klog.Info("Clearing vSphere TagManager cache")

		tm.ensuredTagCategoriesIDs = make(map[string]string)
		tm.ensuredTagIDsForCategoryName = make(map[string]map[string]string)
		tm.regionTagNameToTag = make(map[string]*tags.Tag)
		tm.zoneTagNameToTag = make(map[string]*tags.Tag)
		tm.tagIdToDcRef = make(map[string]object.Reference)
		tm.tagIdToClusterRef = make(map[string]object.Reference)
	}, wait.Jitter(6*time.Hour, 0))

	return tm
}

func (tm *TagManager) GetObjectRefsByTag(ctx context.Context, client *govmomi.Client, restClient *rest.Client, tagName string) ([]object.Reference, error) {
	var objects []mo.Reference
	var objectRefs []object.Reference

	tagsClient := tags.NewManager(restClient)

	var err error
	objects, err = tagsClient.ListAttachedObjects(ctx, tagName)
	if err != nil {
		return nil, err
	}

	finder := find.NewFinder(client.Client)

	for _, MO := range objects {
		ref, err := finder.ObjectReference(ctx, MO.Reference())
		if err != nil {
			return nil, err
		}

		objectRefs = append(objectRefs, ref)
	}

	return objectRefs, nil
}

func (tm *TagManager) EnsureAttachedTagToObject(ctx context.Context, restClient *rest.Client, tagName, categoryName string, ref mo.Reference) error {
	tagsClient := tags.NewManager(restClient)

	if tagID, ok := tm.ensuredTagIDsForCategoryName[categoryName][tagName]; ok {
		attachedTags, err := tagsClient.GetAttachedTags(ctx, ref)

		var tagIsAttached bool
		for _, attachedTag := range attachedTags {
			tagIsAttached = attachedTag.ID == tagID
		}
		if tagIsAttached {
			return nil
		}

		err = tagsClient.AttachTag(ctx, tagID, ref)
		if err != nil {
			return err
		}

		return nil
	}

	tagToAttach, err := tagsClient.GetTagForCategory(ctx, tagName, categoryName)
	if err != nil {
		return err
	}

	if _, ok := tm.ensuredTagIDsForCategoryName[categoryName]; !ok {
		tm.ensuredTagIDsForCategoryName[categoryName] = make(map[string]string)
	}
	tm.ensuredTagIDsForCategoryName[categoryName][tagName] = tagToAttach.ID

	attachedTags, err := tagsClient.GetAttachedTags(ctx, ref)

	var tagIsAttached bool
	for _, attachedTag := range attachedTags {
		tagIsAttached = attachedTag.Name == tagName
	}
	if tagIsAttached {
		return nil
	}

	err = tagsClient.AttachTag(ctx, tagToAttach.ID, ref)
	if err != nil {
		return err
	}

	return nil
}

func (tm *TagManager) EnsureTagCategory(ctx context.Context, restClient *rest.Client, categoryName, description, cardinality string, associableTypes []string) error {
	tm.Lock()
	defer tm.Unlock()

	if _, ok := tm.ensuredTagCategoriesIDs[categoryName]; ok {
		return nil
	}

	tagsClient := tags.NewManager(restClient)

	tagCategory, err := tagsClient.GetCategory(ctx, categoryName)
	if err == nil {
		tm.ensuredTagCategoriesIDs[categoryName] = tagCategory.ID
		return nil
	}

	newTagCategoryID, err := tagsClient.CreateCategory(ctx, &tags.Category{
		Name:            categoryName,
		Description:     description,
		Cardinality:     cardinality,
		AssociableTypes: associableTypes,
	})
	if err != nil {
		return err
	}
	tm.ensuredTagCategoriesIDs[categoryName] = newTagCategoryID

	return nil
}

func (tm *TagManager) EnsureTag(ctx context.Context, restClient *rest.Client, tagName, categoryName, description string) error {
	tm.Lock()
	defer tm.Unlock()

	if _, ok := tm.ensuredTagIDsForCategoryName[categoryName][tagName]; ok {
		return nil
	}

	tagsClient := tags.NewManager(restClient)

	tag, err := tagsClient.GetTagForCategory(ctx, tagName, categoryName)
	if err == nil {
		if _, ok := tm.ensuredTagIDsForCategoryName[categoryName]; !ok {
			tm.ensuredTagIDsForCategoryName[categoryName] = make(map[string]string)
		}
		tm.ensuredTagIDsForCategoryName[categoryName][tagName] = tag.ID
		return nil
	}

	category, err := tagsClient.GetCategory(ctx, categoryName)
	if err != nil {
		return errors.Wrapf(err, "can't get category for categoryName %q", categoryName)
	}

	newTagId, err := tagsClient.CreateTag(ctx, &tags.Tag{
		Description: description,
		Name:        tagName,
		CategoryID:  category.ID,
	})
	if err != nil {
		return errors.Wrapf(err, "can't create tag for categoryName %q", categoryName)
	}

	if _, ok := tm.ensuredTagIDsForCategoryName[categoryName]; !ok {
		tm.ensuredTagIDsForCategoryName[categoryName] = make(map[string]string)
	}
	tm.ensuredTagIDsForCategoryName[categoryName][tagName] = newTagId

	return nil
}

func (tm *TagManager) GetDcByRegion(ctx context.Context, client *govmomi.Client, restClient *rest.Client, region string) (*object.Datacenter, error) {
	tm.Lock()
	defer tm.Unlock()

	var (
		tagsClient = tags.NewManager(restClient)
		regionTag  *tags.Tag

		err error
	)

	if tag, ok := tm.regionTagNameToTag[region]; ok {
		regionTag = tag
	} else {
		regionTag, err = tagsClient.GetTagForCategory(ctx, region, tm.regionTagCategoryName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get region tag for category")
		}
		tm.regionTagNameToTag[region] = regionTag
	}

	finder := find.NewFinder(client.Client)

	if dc, ok := tm.tagIdToDcRef[regionTag.ID]; ok {
		datacenterRef, err := finder.ObjectReference(ctx, dc.Reference())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get DataCenter by datacenterRef: %+v", dc.Reference())
		}

		datacenter := datacenterRef.(*object.Datacenter)
		return datacenter, nil
	}

	taggedObjects, err := tagsClient.ListAttachedObjects(ctx, regionTag.ID)
	var dcs []*object.Datacenter
	for _, taggedObject := range taggedObjects {
		if taggedObject.Reference().Type == "Datacenter" {
			dcRef, err := finder.ObjectReference(ctx, taggedObject.Reference())
			if err != nil {
				return nil, err
			}

			dcs = append(dcs, dcRef.(*object.Datacenter))
		}
	}
	if len(dcs) != 1 {
		return nil, fmt.Errorf("only one DC should match \"region\" tag. Tag: \"%s\", DCs: \"%v\"", regionTag.Name, func() []string {
			datacenters := make([]string, 0, len(dcs))
			for _, dc := range dcs {
				datacenters = append(datacenters, object.NewDatacenter(client.Client, dc.Reference()).InventoryPath)
			}
			return datacenters
		}())
	}

	tm.tagIdToDcRef[regionTag.ID] = dcs[0].Reference()

	return dcs[0], nil
}

func (tm *TagManager) GetClusterByZone(ctx context.Context, client *govmomi.Client, restClient *rest.Client, zone string) (*object.ClusterComputeResource, error) {
	tm.Lock()
	defer tm.Unlock()

	var (
		tagsClient = tags.NewManager(restClient)
		zoneTag    *tags.Tag

		err error
	)

	if tag, ok := tm.zoneTagNameToTag[zone]; ok {
		zoneTag = tag
	} else {
		zoneTag, err = tagsClient.GetTagForCategory(ctx, zone, tm.zoneTagCategoryName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get zone tag for category")
		}
		tm.zoneTagNameToTag[zone] = zoneTag
	}

	finder := find.NewFinder(client.Client)

	if cluster, ok := tm.tagIdToClusterRef[zoneTag.ID]; ok {
		clsRef, err := finder.ObjectReference(ctx, cluster.Reference())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get Cluster by clusterRef: %+v", cluster.Reference())
		}

		cls := clsRef.(*object.ClusterComputeResource)
		return cls, nil
	}

	taggedObjects, err := tagsClient.ListAttachedObjects(ctx, zoneTag.ID)

	var clusters []*object.ClusterComputeResource
	for _, taggedObject := range taggedObjects {
		if taggedObject.Reference().Type == "ClusterComputeResource" {
			clusterRef, err := finder.ObjectReference(ctx, taggedObject.Reference())
			if err != nil {
				return nil, err
			}

			clusters = append(clusters, clusterRef.(*object.ClusterComputeResource))
		}
	}

	if len(clusters) != 1 {
		return nil, fmt.Errorf("only one Cluster should match \"zone\" tag. Tag: \"%s\", Clusters: \"%v\"", zoneTag.Name, func() []string {
			cs := make([]string, 0, len(clusters))
			for _, c := range clusters {
				cs = append(cs, object.NewDatacenter(client.Client, c.Reference()).InventoryPath)
			}
			return cs
		}())
	}

	clusterRef, err := finder.ObjectReference(ctx, clusters[0].Reference())
	if err != nil {
		return nil, err
	}

	cluster := clusterRef.(*object.ClusterComputeResource)
	tm.tagIdToClusterRef[zoneTag.ID] = cluster.Reference()

	return cluster, nil
}
