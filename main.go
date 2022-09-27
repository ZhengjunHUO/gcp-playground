package main

import (
	"fmt"
	"context"

	"google.golang.org/api/compute/v1"

	"github.com/ZhengjunHUO/gcp-playground/pkg"
)

const (
	ZERO int64 = 0
	ONE  int64 = 1
)

var (
	PROJECT_NAME	= "ENTER_YOUR_PROJECT_NAME_HERE"
	labelKey	= "foo"
	labelVal	= "bar"
)

func main() {
	gkecluster := &pkg.GKECluster{ProjectName: PROJECT_NAME, Ctx: context.Background(),}

	gkecluster.FindCluster(labelKey, labelVal)
	if gkecluster.Cluster == nil {
		fmt.Printf("No cluster matches the label [%s=%s] ! Quit...\n", labelKey, labelVal)
		return
	}

	instGrps := gkecluster.ListInstanceGroups()


	// Get a compute service
	ctx := context.Background()
	cmpsvc, err := compute.NewService(ctx)
	if err != nil {
		fmt.Println("compute NewService error: ", err)
		return
	}

	// Get an InstanceGroupManagersService
	grpMgrSVC := compute.NewInstanceGroupManagersService(cmpsvc)

	for _, v := range instGrps {
		// Get InstanceGroupManager
		mgr, err := grpMgrSVC.Get(v.Projet, v.Zone, v.Manager).Do()
		if err != nil {
			fmt.Println("Get instance group manager error: ", err)
			continue
		}

		fmt.Printf("Found %s/%s/%s with size %d. Status: ", v.Projet, v.Zone, v.Manager, mgr.TargetSize)
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
