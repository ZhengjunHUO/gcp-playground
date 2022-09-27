package main

import (
	"fmt"
	"context"
	"strings"

	"google.golang.org/api/container/v1"
	"google.golang.org/api/compute/v1"
)

const (
	PROJECT_NAME = "ENTER_YOUR_PROJECT_NAME_HERE"
	SEPERATE_STR = "zones"
	ZERO   int64 = 0
	ONE    int64 = 1
)

var (
	labelKey = "foo"
	labelVal = "bar"
)

type instanceGroup struct {
	projet  string
	zone    string
	manager string
}

func main() {
	ctx := context.Background()

	// Get a container service
	ctnsvc, err := container.NewService(ctx)
	if err != nil {
		fmt.Println("container NewService error: ", err)
		return
	}

	// Get a cluster service
	clusterSVC := container.NewProjectsZonesClustersService(ctnsvc)
	// Get a ListClustersResponse contains list of cluster objects in ALL zones with detail under PROJECT_NAME
	resp, err := clusterSVC.List(PROJECT_NAME, "-").Do()
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
	var targetCluster *container.Cluster
	for _, cluster := range resp.Clusters {
		if v, ok := cluster.ResourceLabels[labelKey]; ok && v == labelVal {
			targetCluster = cluster
			break
		}
	}

	if targetCluster == nil {
		fmt.Printf("No cluster matches the label [%s=%s] ! Quit...\n", labelKey, labelVal)
		return
	}
	/* Print selected cluster's detail
	buf, err := targetCluster.MarshalJSON()
	if err != nil {
		fmt.Println("Cluster Marshall error: ", err)
		return
	}

	fmt.Println(string(buf))
	*/

	// Find a list of instance group associated to the cluster
	instGrps := []instanceGroup{}
	for _, nodepool := range targetCluster.NodePools {
		for _, str := range nodepool.InstanceGroupUrls {
			// need Go 1.18
			//_, after, _ := strings.Cut(str, SEPERATE_STR)
			strs := strings.Split(str[strings.Index(str, SEPERATE_STR)+len(SEPERATE_STR)+1:], "/")
			if len(strs) < 3 {
				fmt.Printf("Can't get instance groups' info correctly, got: %v\n", strs)
				continue
			}
			instGrps = append(instGrps, instanceGroup{PROJECT_NAME, strs[0], strs[2]})
		}
	}
	/* Deprecated
	for _, str := range targetCluster.InstanceGroupUrls {
		strs := strings.Split(str[strings.Index(str, SEPERATE_STR)+len(SEPERATE_STR)+1:], "/")
		if len(strs) < 3 {
			fmt.Printf("Can't get instance groups' info correctly, got: %v\n", strs)
			continue
		}
		instGrps = append(instGrps, instanceGroup{PROJECT_NAME, strs[0], strs[2]})
	}
	*/

	// Get a compute service
	cmpsvc, err := compute.NewService(ctx)
	if err != nil {
		fmt.Println("compute NewService error: ", err)
		return
	}

	// Get an InstanceGroupManagersService
	grpMgrSVC := compute.NewInstanceGroupManagersService(cmpsvc)

	for _, v := range instGrps {
		// Get InstanceGroupManager
		mgr, err := grpMgrSVC.Get(v.projet, v.zone, v.manager).Do()
		if err != nil {
			fmt.Println("Get instance group manager error: ", err)
			continue
		}

		fmt.Printf("Found %s/%s/%s with size %d. Status: ", v.projet, v.zone, v.manager, mgr.TargetSize)
		if mgr.Status.IsStable {
			fmt.Println("OK")
		}else{
			fmt.Println("KO")
		}

		/* Print InstanceGroupManager's detail
		buf, err := mgr.MarshalJSON()
		if err != nil {
			fmt.Println("Marshall error: ", err)
			continue
		}

		fmt.Println(string(buf))
		fmt.Println("==================\n")
		*/
	}

	/* resize group manager to zero
	for _, chosen := range instGrps {
		fmt.Printf("Resize %s/%s/%s to 0\n", chosen.projet, chosen.zone, chosen.manager)
		_, err := grpMgrSVC.Resize(chosen.projet, chosen.zone, chosen.manager, ZERO).Do()
		if err != nil {
			fmt.Println("Resize to zero err: ", err)
		}
	}
	*/
}
