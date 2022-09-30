package pkg

import (
	"fmt"
	"strings"

	"google.golang.org/api/container/v1"
)

const (
	SEPERATE_STR = "zones"
)

func (g *GKECluster) FindCluster(labelKey, labelVal string) {
	// Get a container service
	ctnsvc, err := container.NewService(g.Ctx)
	if err != nil {
		fmt.Println("container NewService error: ", err)
		return
	}

	// Get a cluster service
	clusterSVC := container.NewProjectsZonesClustersService(ctnsvc)
	// Get a ListClustersResponse contains list of cluster objects in ALL zones with detail under PROJECT_NAME
	resp, err := clusterSVC.List(g.ProjectName, "-").Do()
	if err != nil {
		fmt.Println("List cluster error: ", err)
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
		fmt.Println("[WARNING] Wait for a non-nil cluster !")
		return instGrps
	}

	for _, nodepool := range g.Cluster.NodePools {
		for _, str := range nodepool.InstanceGroupUrls {
			// need Go 1.18
			//_, after, _ := strings.Cut(str, SEPERATE_STR)
			strs := strings.Split(str[strings.Index(str, SEPERATE_STR)+len(SEPERATE_STR)+1:], "/")
			if len(strs) < 3 {
				fmt.Printf("Can't get instance groups' info correctly, got: %v\n", strs)
				continue
			}
			instGrps = append(instGrps, &InstanceGroup{g.ProjectName, strs[0], strs[2], nil})
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
