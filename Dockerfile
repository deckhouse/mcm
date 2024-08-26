#############      builder                                  #############
FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5@sha256:5e9360b0462b3717384bbb4330bc75b53dc9c11dabdf39c5834fe8708439f6d2 AS builder

WORKDIR /go/src/github.com/gardener/machine-controller-manager
COPY . .

RUN .ci/build

#############      base                                     #############
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.1@sha256:d7342993700f8cd7aba8496c2d0e57be0666e80b4c441925fc6f9361fa81d10e as base

RUN apk add --update bash curl tzdata
WORKDIR /

#############      machine-controller-manager               #############
FROM base AS machine-controller-manager

COPY --from=builder /go/src/github.com/gardener/machine-controller-manager/bin/rel/machine-controller-manager /machine-controller-manager
ENTRYPOINT ["/machine-controller-manager"]
