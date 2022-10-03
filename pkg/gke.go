package pkg

import (
	"fmt"
	"strings"
	"context"

	"google.golang.org/api/container/v1"
	"google.golang.org/api/compute/v1"
)

const (
	SEPERATE_STR = "zones"
)

func NewGKECluster(projectName string, ctx context.Context) *GKECluster {
	cluster := &GKECluster{
                ProjectName: projectName,
                Ctx:         ctx,
        }

	// Get a container service
	ctnsvc, err := container.NewService(ctx)
	if err != nil {
		fmt.Println("[WARN] container NewService error: ", err)
		return cluster
	}

	cluster.CtnService = ctnsvc
	return cluster
}

func (g *GKECluster) WithComputeService(ctx context.Context) *GKECluster {
        // Get a compute service
        cmpsvc, err := compute.NewService(ctx)
        if err != nil {
                fmt.Println("[WARN] compute NewService error: ", err)
                return g
        }

	g.VmService = cmpsvc
	return g
}

func (g *GKECluster) FindCluster(labelKey, labelVal string) {
	if g.CtnService == nil {
                fmt.Println("[WARN] container Service not initialized properly, skip FindCluster.")
		return
	}

	// Get a cluster service
	clusterSVC := container.NewProjectsZonesClustersService(g.CtnService)
	// Get a ListClustersResponse contains list of cluster objects in ALL zones with detail under PROJECT_NAME
	resp, err := clusterSVC.List(g.ProjectName, "-").Do()
	if err != nil {
		fmt.Println("[WARN] List cluster error: ", err)
		return
	}

	/* Print all clusters' detail
	buf, err := resp.MarshalJSON()
	if err != nil {
		fmt.Println("Marshall error: ", err)
		return
	}

	fmt.Println(string(buf))
	*/

	// Find the Cluster with specific label
	for _, cluster := range resp.Clusters {
		if v, ok := cluster.ResourceLabels[labelKey]; ok && v == labelVal {
			g.Cluster = cluster
			break
		}
	}

	/* Print selected cluster's detail
	buf, err := targetCluster.MarshalJSON()
	if err != nil {
		fmt.Println("Cluster Marshall error: ", err)
		return
	}

	fmt.Println(string(buf))
	*/
}

func (g *GKECluster) ListInstanceGroups() []*InstanceGroup {
	instGrps := []*InstanceGroup{}

	if g.Cluster == nil {
		fmt.Println("[WARN] Wait for a non-nil cluster !")
		return instGrps
	}

	// Get an InstanceGroupManagersService
	grpMgrSVC := compute.NewInstanceGroupManagersService(g.VmService)

	for _, nodepool := range g.Cluster.NodePools {
		for _, str := range nodepool.InstanceGroupUrls {
			// need Go 1.18
			//_, after, _ := strings.Cut(str, SEPERATE_STR)
			strs := strings.Split(str[strings.Index(str, SEPERATE_STR)+len(SEPERATE_STR)+1:], "/")
			if len(strs) < 3 {
				fmt.Printf("[WARN] Can't get instance groups' info correctly, got: %v\n", strs)
				continue
			}

			// Get InstanceGroupManager
			mgr, err := grpMgrSVC.Get(g.ProjectName, strs[0], strs[2]).Do()
			if err != nil {
				fmt.Println("[WARN] Get instance group manager error: ", err)
			}

			instGrps = append(instGrps, &InstanceGroup{g.ProjectName, strs[0], strs[2], mgr, grpMgrSVC})
		}
	}
	/* Deprecated
	for _, str := range targetCluster.InstanceGroupUrls {
		strs := strings.Split(str[strings.Index(str, SEPERATE_STR)+len(SEPERATE_STR)+1:], "/")
		if len(strs) < 3 {
			fmt.Printf("Can't get instance groups' info correctly, got: %v\n", strs)
			continue
		}
		instGrps = append(instGrps, &instanceGroup{PROJECT_NAME, strs[0], strs[2]})
	}
	*/

	return instGrps
}
