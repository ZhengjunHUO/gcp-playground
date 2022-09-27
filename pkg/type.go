package pkg

import (
	"context"
	"google.golang.org/api/container/v1"

)

type GKECluster struct {
	ProjectName	string
	Ctx		context.Context
	Cluster		*container.Cluster
}

type InstanceGroup struct {
	Projet  string
	Zone    string
	Manager string
}

