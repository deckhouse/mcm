// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/mdb/postgresql/v1/user.proto

package postgresql

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	_ "github.com/yandex-cloud/go-genproto/yandex/cloud"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type UserSettings_SynchronousCommit int32

const (
	UserSettings_SYNCHRONOUS_COMMIT_UNSPECIFIED  UserSettings_SynchronousCommit = 0
	UserSettings_SYNCHRONOUS_COMMIT_ON           UserSettings_SynchronousCommit = 1
	UserSettings_SYNCHRONOUS_COMMIT_OFF          UserSettings_SynchronousCommit = 2
	UserSettings_SYNCHRONOUS_COMMIT_LOCAL        UserSettings_SynchronousCommit = 3
	UserSettings_SYNCHRONOUS_COMMIT_REMOTE_WRITE UserSettings_SynchronousCommit = 4
	UserSettings_SYNCHRONOUS_COMMIT_REMOTE_APPLY UserSettings_SynchronousCommit = 5
)

var UserSettings_SynchronousCommit_name = map[int32]string{
	0: "SYNCHRONOUS_COMMIT_UNSPECIFIED",
	1: "SYNCHRONOUS_COMMIT_ON",
	2: "SYNCHRONOUS_COMMIT_OFF",
	3: "SYNCHRONOUS_COMMIT_LOCAL",
	4: "SYNCHRONOUS_COMMIT_REMOTE_WRITE",
	5: "SYNCHRONOUS_COMMIT_REMOTE_APPLY",
}

var UserSettings_SynchronousCommit_value = map[string]int32{
	"SYNCHRONOUS_COMMIT_UNSPECIFIED":  0,
	"SYNCHRONOUS_COMMIT_ON":           1,
	"SYNCHRONOUS_COMMIT_OFF":          2,
	"SYNCHRONOUS_COMMIT_LOCAL":        3,
	"SYNCHRONOUS_COMMIT_REMOTE_WRITE": 4,
	"SYNCHRONOUS_COMMIT_REMOTE_APPLY": 5,
}

func (x UserSettings_SynchronousCommit) String() string {
	return proto.EnumName(UserSettings_SynchronousCommit_name, int32(x))
}

func (UserSettings_SynchronousCommit) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{3, 0}
}

type UserSettings_LogStatement int32

const (
	UserSettings_LOG_STATEMENT_UNSPECIFIED UserSettings_LogStatement = 0
	UserSettings_LOG_STATEMENT_NONE        UserSettings_LogStatement = 1
	UserSettings_LOG_STATEMENT_DDL         UserSettings_LogStatement = 2
	UserSettings_LOG_STATEMENT_MOD         UserSettings_LogStatement = 3
	UserSettings_LOG_STATEMENT_ALL         UserSettings_LogStatement = 4
)

var UserSettings_LogStatement_name = map[int32]string{
	0: "LOG_STATEMENT_UNSPECIFIED",
	1: "LOG_STATEMENT_NONE",
	2: "LOG_STATEMENT_DDL",
	3: "LOG_STATEMENT_MOD",
	4: "LOG_STATEMENT_ALL",
}

var UserSettings_LogStatement_value = map[string]int32{
	"LOG_STATEMENT_UNSPECIFIED": 0,
	"LOG_STATEMENT_NONE":        1,
	"LOG_STATEMENT_DDL":         2,
	"LOG_STATEMENT_MOD":         3,
	"LOG_STATEMENT_ALL":         4,
}

func (x UserSettings_LogStatement) String() string {
	return proto.EnumName(UserSettings_LogStatement_name, int32(x))
}

func (UserSettings_LogStatement) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{3, 1}
}

type UserSettings_TransactionIsolation int32

const (
	UserSettings_TRANSACTION_ISOLATION_UNSPECIFIED      UserSettings_TransactionIsolation = 0
	UserSettings_TRANSACTION_ISOLATION_READ_UNCOMMITTED UserSettings_TransactionIsolation = 1
	UserSettings_TRANSACTION_ISOLATION_READ_COMMITTED   UserSettings_TransactionIsolation = 2
	UserSettings_TRANSACTION_ISOLATION_REPEATABLE_READ  UserSettings_TransactionIsolation = 3
	UserSettings_TRANSACTION_ISOLATION_SERIALIZABLE     UserSettings_TransactionIsolation = 4
)

var UserSettings_TransactionIsolation_name = map[int32]string{
	0: "TRANSACTION_ISOLATION_UNSPECIFIED",
	1: "TRANSACTION_ISOLATION_READ_UNCOMMITTED",
	2: "TRANSACTION_ISOLATION_READ_COMMITTED",
	3: "TRANSACTION_ISOLATION_REPEATABLE_READ",
	4: "TRANSACTION_ISOLATION_SERIALIZABLE",
}

var UserSettings_TransactionIsolation_value = map[string]int32{
	"TRANSACTION_ISOLATION_UNSPECIFIED":      0,
	"TRANSACTION_ISOLATION_READ_UNCOMMITTED": 1,
	"TRANSACTION_ISOLATION_READ_COMMITTED":   2,
	"TRANSACTION_ISOLATION_REPEATABLE_READ":  3,
	"TRANSACTION_ISOLATION_SERIALIZABLE":     4,
}

func (x UserSettings_TransactionIsolation) String() string {
	return proto.EnumName(UserSettings_TransactionIsolation_name, int32(x))
}

func (UserSettings_TransactionIsolation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{3, 2}
}

// A PostgreSQL User resource. For more information, see
// the [Developer's Guide](/docs/managed-postgresql/concepts).
type User struct {
	// Name of the PostgreSQL user.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// ID of the PostgreSQL cluster the user belongs to.
	ClusterId string `protobuf:"bytes,2,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	// Set of permissions granted to the user.
	Permissions []*Permission `protobuf:"bytes,3,rep,name=permissions,proto3" json:"permissions,omitempty"`
	// Number of database connections available to the user.
	ConnLimit int64 `protobuf:"varint,4,opt,name=conn_limit,json=connLimit,proto3" json:"conn_limit,omitempty"`
	// Postgresql settings for this user
	Settings *UserSettings `protobuf:"bytes,5,opt,name=settings,proto3" json:"settings,omitempty"`
	// User can login (default True)
	Login *wrappers.BoolValue `protobuf:"bytes,6,opt,name=login,proto3" json:"login,omitempty"`
	// User grants (GRANT <role> TO <user>), role must be other user
	Grants               []string `protobuf:"bytes,7,rep,name=grants,proto3" json:"grants,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetClusterId() string {
	if m != nil {
		return m.ClusterId
	}
	return ""
}

func (m *User) GetPermissions() []*Permission {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func (m *User) GetConnLimit() int64 {
	if m != nil {
		return m.ConnLimit
	}
	return 0
}

func (m *User) GetSettings() *UserSettings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *User) GetLogin() *wrappers.BoolValue {
	if m != nil {
		return m.Login
	}
	return nil
}

func (m *User) GetGrants() []string {
	if m != nil {
		return m.Grants
	}
	return nil
}

type Permission struct {
	// Name of the database that the permission grants access to.
	DatabaseName         string   `protobuf:"bytes,1,opt,name=database_name,json=databaseName,proto3" json:"database_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Permission) Reset()         { *m = Permission{} }
func (m *Permission) String() string { return proto.CompactTextString(m) }
func (*Permission) ProtoMessage()    {}
func (*Permission) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{1}
}

func (m *Permission) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Permission.Unmarshal(m, b)
}
func (m *Permission) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Permission.Marshal(b, m, deterministic)
}
func (m *Permission) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Permission.Merge(m, src)
}
func (m *Permission) XXX_Size() int {
	return xxx_messageInfo_Permission.Size(m)
}
func (m *Permission) XXX_DiscardUnknown() {
	xxx_messageInfo_Permission.DiscardUnknown(m)
}

var xxx_messageInfo_Permission proto.InternalMessageInfo

func (m *Permission) GetDatabaseName() string {
	if m != nil {
		return m.DatabaseName
	}
	return ""
}

type UserSpec struct {
	// Name of the PostgreSQL user.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Password of the PostgreSQL user.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// Set of permissions to grant to the user.
	Permissions []*Permission `protobuf:"bytes,3,rep,name=permissions,proto3" json:"permissions,omitempty"`
	// Number of database connections that should be available to the user.
	ConnLimit *wrappers.Int64Value `protobuf:"bytes,4,opt,name=conn_limit,json=connLimit,proto3" json:"conn_limit,omitempty"`
	// Postgresql settings for this user
	Settings *UserSettings `protobuf:"bytes,5,opt,name=settings,proto3" json:"settings,omitempty"`
	// User can login (default True)
	Login *wrappers.BoolValue `protobuf:"bytes,6,opt,name=login,proto3" json:"login,omitempty"`
	// User grants (GRANT <role> TO <user>), role must be other user
	Grants               []string `protobuf:"bytes,7,rep,name=grants,proto3" json:"grants,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserSpec) Reset()         { *m = UserSpec{} }
func (m *UserSpec) String() string { return proto.CompactTextString(m) }
func (*UserSpec) ProtoMessage()    {}
func (*UserSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{2}
}

func (m *UserSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserSpec.Unmarshal(m, b)
}
func (m *UserSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserSpec.Marshal(b, m, deterministic)
}
func (m *UserSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserSpec.Merge(m, src)
}
func (m *UserSpec) XXX_Size() int {
	return xxx_messageInfo_UserSpec.Size(m)
}
func (m *UserSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_UserSpec.DiscardUnknown(m)
}

var xxx_messageInfo_UserSpec proto.InternalMessageInfo

func (m *UserSpec) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UserSpec) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *UserSpec) GetPermissions() []*Permission {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func (m *UserSpec) GetConnLimit() *wrappers.Int64Value {
	if m != nil {
		return m.ConnLimit
	}
	return nil
}

func (m *UserSpec) GetSettings() *UserSettings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *UserSpec) GetLogin() *wrappers.BoolValue {
	if m != nil {
		return m.Login
	}
	return nil
}

func (m *UserSpec) GetGrants() []string {
	if m != nil {
		return m.Grants
	}
	return nil
}

// Postgresql user settings config
type UserSettings struct {
	DefaultTransactionIsolation UserSettings_TransactionIsolation `protobuf:"varint,1,opt,name=default_transaction_isolation,json=defaultTransactionIsolation,proto3,enum=yandex.cloud.mdb.postgresql.v1.UserSettings_TransactionIsolation" json:"default_transaction_isolation,omitempty"`
	// in milliseconds.
	LockTimeout *wrappers.Int64Value `protobuf:"bytes,2,opt,name=lock_timeout,json=lockTimeout,proto3" json:"lock_timeout,omitempty"`
	// in milliseconds.
	LogMinDurationStatement *wrappers.Int64Value           `protobuf:"bytes,3,opt,name=log_min_duration_statement,json=logMinDurationStatement,proto3" json:"log_min_duration_statement,omitempty"`
	SynchronousCommit       UserSettings_SynchronousCommit `protobuf:"varint,4,opt,name=synchronous_commit,json=synchronousCommit,proto3,enum=yandex.cloud.mdb.postgresql.v1.UserSettings_SynchronousCommit" json:"synchronous_commit,omitempty"`
	// in bytes.
	TempFileLimit        *wrappers.Int64Value      `protobuf:"bytes,5,opt,name=temp_file_limit,json=tempFileLimit,proto3" json:"temp_file_limit,omitempty"`
	LogStatement         UserSettings_LogStatement `protobuf:"varint,6,opt,name=log_statement,json=logStatement,proto3,enum=yandex.cloud.mdb.postgresql.v1.UserSettings_LogStatement" json:"log_statement,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *UserSettings) Reset()         { *m = UserSettings{} }
func (m *UserSettings) String() string { return proto.CompactTextString(m) }
func (*UserSettings) ProtoMessage()    {}
func (*UserSettings) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddb893fd4eb810a7, []int{3}
}

func (m *UserSettings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserSettings.Unmarshal(m, b)
}
func (m *UserSettings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserSettings.Marshal(b, m, deterministic)
}
func (m *UserSettings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserSettings.Merge(m, src)
}
func (m *UserSettings) XXX_Size() int {
	return xxx_messageInfo_UserSettings.Size(m)
}
func (m *UserSettings) XXX_DiscardUnknown() {
	xxx_messageInfo_UserSettings.DiscardUnknown(m)
}

var xxx_messageInfo_UserSettings proto.InternalMessageInfo

func (m *UserSettings) GetDefaultTransactionIsolation() UserSettings_TransactionIsolation {
	if m != nil {
		return m.DefaultTransactionIsolation
	}
	return UserSettings_TRANSACTION_ISOLATION_UNSPECIFIED
}

func (m *UserSettings) GetLockTimeout() *wrappers.Int64Value {
	if m != nil {
		return m.LockTimeout
	}
	return nil
}

func (m *UserSettings) GetLogMinDurationStatement() *wrappers.Int64Value {
	if m != nil {
		return m.LogMinDurationStatement
	}
	return nil
}

func (m *UserSettings) GetSynchronousCommit() UserSettings_SynchronousCommit {
	if m != nil {
		return m.SynchronousCommit
	}
	return UserSettings_SYNCHRONOUS_COMMIT_UNSPECIFIED
}

func (m *UserSettings) GetTempFileLimit() *wrappers.Int64Value {
	if m != nil {
		return m.TempFileLimit
	}
	return nil
}

func (m *UserSettings) GetLogStatement() UserSettings_LogStatement {
	if m != nil {
		return m.LogStatement
	}
	return UserSettings_LOG_STATEMENT_UNSPECIFIED
}

func init() {
	proto.RegisterEnum("yandex.cloud.mdb.postgresql.v1.UserSettings_SynchronousCommit", UserSettings_SynchronousCommit_name, UserSettings_SynchronousCommit_value)
	proto.RegisterEnum("yandex.cloud.mdb.postgresql.v1.UserSettings_LogStatement", UserSettings_LogStatement_name, UserSettings_LogStatement_value)
	proto.RegisterEnum("yandex.cloud.mdb.postgresql.v1.UserSettings_TransactionIsolation", UserSettings_TransactionIsolation_name, UserSettings_TransactionIsolation_value)
	proto.RegisterType((*User)(nil), "yandex.cloud.mdb.postgresql.v1.User")
	proto.RegisterType((*Permission)(nil), "yandex.cloud.mdb.postgresql.v1.Permission")
	proto.RegisterType((*UserSpec)(nil), "yandex.cloud.mdb.postgresql.v1.UserSpec")
	proto.RegisterType((*UserSettings)(nil), "yandex.cloud.mdb.postgresql.v1.UserSettings")
}

func init() {
	proto.RegisterFile("yandex/cloud/mdb/postgresql/v1/user.proto", fileDescriptor_ddb893fd4eb810a7)
}

var fileDescriptor_ddb893fd4eb810a7 = []byte{
	// 924 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x56, 0x5f, 0x6f, 0xdb, 0x54,
	0x14, 0xc7, 0x4d, 0x5a, 0x9a, 0x93, 0x64, 0xb8, 0x57, 0x6c, 0xa4, 0x19, 0x29, 0x21, 0x63, 0x53,
	0x5a, 0x11, 0xa7, 0xce, 0xd0, 0xb4, 0x09, 0x56, 0xc9, 0x49, 0x5c, 0x66, 0xc9, 0xb1, 0x23, 0xdb,
	0x05, 0x56, 0x04, 0x96, 0x13, 0xdf, 0x7a, 0x16, 0xb6, 0x6f, 0xf0, 0xb5, 0x3b, 0xc6, 0x3b, 0x2f,
	0x7d, 0xe5, 0x63, 0xf0, 0x41, 0xba, 0x2f, 0x00, 0x1f, 0x01, 0xf1, 0xcc, 0x23, 0x4f, 0x28, 0x76,
	0xd2, 0x24, 0x4b, 0x96, 0xaa, 0x12, 0x2f, 0xbc, 0xd9, 0xe7, 0xfc, 0x7e, 0xe7, 0xcf, 0xef, 0x9c,
	0x23, 0x5d, 0xd8, 0x7f, 0x65, 0x05, 0x36, 0xfe, 0xa9, 0x39, 0xf4, 0x48, 0x6c, 0x37, 0x7d, 0x7b,
	0xd0, 0x1c, 0x11, 0x1a, 0x39, 0x21, 0xa6, 0x3f, 0x7a, 0xcd, 0x73, 0xbe, 0x19, 0x53, 0x1c, 0x72,
	0xa3, 0x90, 0x44, 0x04, 0xed, 0xa5, 0x50, 0x2e, 0x81, 0x72, 0xbe, 0x3d, 0xe0, 0x66, 0x50, 0xee,
	0x9c, 0x2f, 0xef, 0x39, 0x84, 0x38, 0x1e, 0x6e, 0x26, 0xe8, 0x41, 0x7c, 0xd6, 0x7c, 0x19, 0x5a,
	0xa3, 0x11, 0x0e, 0x69, 0xca, 0x2f, 0x57, 0x16, 0x52, 0x9d, 0x5b, 0x9e, 0x6b, 0x5b, 0x91, 0x4b,
	0x82, 0xd4, 0x5d, 0xfb, 0x7d, 0x03, 0xb2, 0x27, 0x14, 0x87, 0x08, 0x41, 0x36, 0xb0, 0x7c, 0x5c,
	0x62, 0xaa, 0x4c, 0x3d, 0xa7, 0x25, 0xdf, 0xa8, 0x02, 0x30, 0xf4, 0x62, 0x1a, 0xe1, 0xd0, 0x74,
	0xed, 0xd2, 0x46, 0xe2, 0xc9, 0x4d, 0x2c, 0x92, 0x8d, 0x64, 0xc8, 0x8f, 0x70, 0xe8, 0xbb, 0x94,
	0xba, 0x24, 0xa0, 0xa5, 0x4c, 0x35, 0x53, 0xcf, 0xb7, 0x0e, 0xb8, 0xf5, 0x05, 0x73, 0xfd, 0x2b,
	0x8a, 0x36, 0x4f, 0x4f, 0x92, 0x91, 0x20, 0x30, 0x3d, 0xd7, 0x77, 0xa3, 0x52, 0xb6, 0xca, 0xd4,
	0x33, 0x5a, 0x6e, 0x6c, 0x91, 0xc7, 0x06, 0xf4, 0x0c, 0xb6, 0x29, 0x8e, 0x22, 0x37, 0x70, 0x68,
	0x69, 0xb3, 0xca, 0xd4, 0xf3, 0xad, 0x4f, 0xaf, 0xcb, 0x34, 0xee, 0x4b, 0x9f, 0x70, 0xb4, 0x2b,
	0x36, 0x3a, 0x84, 0x4d, 0x8f, 0x38, 0x6e, 0x50, 0xda, 0x4a, 0xc2, 0x94, 0xb9, 0x54, 0x41, 0x6e,
	0xaa, 0x20, 0xd7, 0x26, 0xc4, 0xfb, 0xca, 0xf2, 0x62, 0xac, 0xa5, 0x40, 0xc4, 0xc3, 0x96, 0x13,
	0x5a, 0x41, 0x44, 0x4b, 0xef, 0x56, 0x33, 0xf5, 0x5c, 0x7b, 0xf7, 0xef, 0x4b, 0xbe, 0xf8, 0xad,
	0xd5, 0xf8, 0x59, 0x68, 0x9c, 0x1e, 0x36, 0x9e, 0x98, 0xdf, 0x1d, 0x5c, 0xbc, 0xe6, 0xb3, 0x5f,
	0x3c, 0x7d, 0xf4, 0x50, 0x9b, 0x00, 0x6b, 0x3c, 0xc0, 0xac, 0x51, 0x74, 0x0f, 0x8a, 0xb6, 0x15,
	0x59, 0x03, 0x8b, 0x62, 0x73, 0x4e, 0xe5, 0xc2, 0xd4, 0xa8, 0x58, 0x3e, 0xae, 0xfd, 0x96, 0x81,
	0xed, 0xa4, 0xe4, 0x11, 0x1e, 0x22, 0x7e, 0x7e, 0x1c, 0xed, 0xca, 0x5f, 0x97, 0x3c, 0xf3, 0xf6,
	0xa4, 0xe9, 0xb4, 0xf6, 0x61, 0x7b, 0x64, 0x51, 0xfa, 0x92, 0x84, 0x93, 0x59, 0xb5, 0x8b, 0x63,
	0xda, 0xc5, 0x6b, 0x7e, 0xf3, 0x71, 0x83, 0x6f, 0x3d, 0xd6, 0xae, 0xdc, 0xff, 0xf1, 0xe4, 0xba,
	0x4b, 0x93, 0xcb, 0xb7, 0xee, 0x2e, 0xa9, 0x2a, 0x05, 0xd1, 0xa3, 0xcf, 0x12, 0x59, 0xdb, 0xdb,
	0xff, 0x5c, 0xf2, 0xd9, 0xa3, 0xa7, 0xfc, 0xe1, 0xff, 0x78, 0xc0, 0xbf, 0xe6, 0xa0, 0x30, 0x9f,
	0x1f, 0xfd, 0xc2, 0x40, 0xc5, 0xc6, 0x67, 0x56, 0xec, 0x45, 0x66, 0x14, 0x5a, 0x01, 0xb5, 0x86,
	0xe3, 0x3b, 0x33, 0x5d, 0x4a, 0xbc, 0xe4, 0xe2, 0x92, 0x59, 0xde, 0x6a, 0x09, 0x37, 0xe9, 0x8a,
	0x33, 0x66, 0x91, 0xa4, 0x69, 0x20, 0xed, 0xee, 0x24, 0xcf, 0x2a, 0x27, 0x3a, 0x82, 0x82, 0x47,
	0x86, 0x3f, 0x98, 0x91, 0xeb, 0x63, 0x12, 0x47, 0xc9, 0x2a, 0xac, 0x9f, 0x87, 0x96, 0x1f, 0x13,
	0x8c, 0x14, 0x8f, 0xbe, 0x81, 0xb2, 0x47, 0x1c, 0xd3, 0x77, 0x03, 0xd3, 0x8e, 0xc3, 0x24, 0xa6,
	0x49, 0x23, 0x2b, 0xc2, 0x3e, 0x0e, 0xa2, 0x52, 0xe6, 0xfa, 0x68, 0x1f, 0x78, 0xc4, 0xe9, 0xb9,
	0x41, 0x77, 0x42, 0xd6, 0xa7, 0x5c, 0xe4, 0x03, 0xa2, 0xaf, 0x82, 0xe1, 0x8b, 0x90, 0x04, 0x24,
	0xa6, 0xe6, 0x90, 0xf8, 0xd3, 0x7d, 0xb9, 0xd5, 0x3a, 0xba, 0x91, 0x2a, 0xfa, 0x2c, 0x4c, 0x27,
	0x89, 0xa2, 0xed, 0xd0, 0x37, 0x4d, 0xa8, 0x03, 0xef, 0x45, 0xd8, 0x1f, 0x99, 0x67, 0xae, 0x87,
	0x27, 0xbb, 0xb9, 0x79, 0x7d, 0xf5, 0xc5, 0x31, 0xe7, 0xd8, 0xf5, 0x70, 0xba, 0x95, 0xdf, 0x43,
	0x71, 0xac, 0xc6, 0x4c, 0x80, 0xad, 0xa4, 0xdc, 0x27, 0x37, 0x2a, 0x57, 0x26, 0xce, 0x95, 0x0a,
	0x5a, 0xc1, 0x9b, 0xfb, 0xab, 0xfd, 0xc1, 0xc0, 0xce, 0x52, 0x37, 0xa8, 0x06, 0x7b, 0xfa, 0x73,
	0xa5, 0xf3, 0x4c, 0x53, 0x15, 0xf5, 0x44, 0x37, 0x3b, 0x6a, 0xaf, 0x27, 0x19, 0xe6, 0x89, 0xa2,
	0xf7, 0xc5, 0x8e, 0x74, 0x2c, 0x89, 0x5d, 0xf6, 0x1d, 0xb4, 0x0b, 0xb7, 0x57, 0x60, 0x54, 0x85,
	0x65, 0x50, 0x19, 0xee, 0xac, 0x72, 0x1d, 0x1f, 0xb3, 0x1b, 0xe8, 0x43, 0x28, 0xad, 0xf0, 0xc9,
	0x6a, 0x47, 0x90, 0xd9, 0x0c, 0xba, 0x07, 0x1f, 0xad, 0xf0, 0x6a, 0x62, 0x4f, 0x35, 0x44, 0xf3,
	0x6b, 0x4d, 0x32, 0x44, 0x36, 0xbb, 0x1e, 0x24, 0xf4, 0xfb, 0xf2, 0x73, 0x76, 0xb3, 0x76, 0xc1,
	0x40, 0x61, 0xbe, 0x6f, 0x54, 0x81, 0x5d, 0x59, 0xfd, 0xd2, 0xd4, 0x0d, 0xc1, 0x10, 0x7b, 0xa2,
	0xf2, 0x66, 0x3b, 0x77, 0x00, 0x2d, 0xba, 0x15, 0x55, 0x11, 0x59, 0x06, 0xdd, 0x86, 0x9d, 0x45,
	0x7b, 0xb7, 0x2b, 0xb3, 0x1b, 0xcb, 0xe6, 0x9e, 0xda, 0x65, 0x33, 0xcb, 0x66, 0x41, 0x96, 0xd9,
	0x6c, 0xed, 0x4f, 0x06, 0xde, 0x5f, 0x79, 0x2c, 0xf7, 0xe1, 0x63, 0x43, 0x13, 0x14, 0x5d, 0xe8,
	0x18, 0x92, 0xaa, 0x98, 0x92, 0xae, 0xca, 0x42, 0xf2, 0xb5, 0x58, 0xdc, 0x01, 0x3c, 0x58, 0x0d,
	0xd3, 0x44, 0xa1, 0x6b, 0x9e, 0x28, 0xa9, 0x04, 0x86, 0xd8, 0x65, 0x19, 0x54, 0x87, 0x4f, 0xd6,
	0x60, 0x67, 0xc8, 0x0d, 0xb4, 0x0f, 0xf7, 0xdf, 0x86, 0xec, 0x8b, 0x82, 0x21, 0xb4, 0x65, 0x31,
	0x21, 0xb1, 0x19, 0xf4, 0x00, 0x6a, 0xab, 0xa1, 0xba, 0xa8, 0x49, 0x82, 0x2c, 0x9d, 0x8e, 0xc1,
	0x6c, 0xb6, 0x4d, 0xa1, 0xb6, 0xb0, 0x98, 0xd6, 0xc8, 0x5d, 0x5e, 0xce, 0xd3, 0x9e, 0xe3, 0x46,
	0x2f, 0xe2, 0x01, 0x37, 0x24, 0x7e, 0x33, 0x85, 0x37, 0xd2, 0xe7, 0x81, 0x43, 0x1a, 0x0e, 0x0e,
	0x92, 0xb3, 0x68, 0xae, 0x7f, 0xa2, 0x7c, 0x3e, 0xfb, 0x1b, 0x6c, 0x25, 0x84, 0x87, 0xff, 0x06,
	0x00, 0x00, 0xff, 0xff, 0x71, 0x07, 0xfd, 0xc4, 0xd6, 0x08, 0x00, 0x00,
}