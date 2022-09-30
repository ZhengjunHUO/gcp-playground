package pkg

import (
	"context"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/compute/v1"
)

type GKECluster struct {
	ProjectName	string
	Ctx		context.Context
	Cluster		*container.Cluster
}

type InstanceGroup struct {
	Project string
	Zone    string
	Manager string
	Igm	*compute.InstanceGroupManager
}

