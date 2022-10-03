package main

import (
	"fmt"
	"context"

	"github.com/ZhengjunHUO/gcp-playground/pkg"
)

const (
	ZERO int64 = 0
	ONE  int64 = 1
)

var (
	PROJECT_NAME	= "opensee-ci"
	labelKey	= "env"
	labelVal	= "k8s-vault"
)

func main() {
	/* Use case 1 */
	ctx := context.Background()
	gkecluster := pkg.NewGKECluster(PROJECT_NAME, ctx).WithComputeService(ctx)

	gkecluster.FindCluster(labelKey, labelVal)
	if gkecluster.Cluster == nil {
		fmt.Printf("No cluster matches the label [%s=%s] ! Quit...\n", labelKey, labelVal)
		return
	}

	//fmt.Println(gkecluster.Cluster.ResourceLabels)
	instGrps := gkecluster.ListInstanceGroups()
	for _, v := range instGrps {
		if v.Igm != nil {
			fmt.Printf("[INFO] Found %s/%s/%s with size %d\n", v.Project, v.Zone, v.Manager, v.Igm.TargetSize)
		}
	}


	/* Use case 2 

	//flt := fmt.Sprintf("(labels.%s=%s)", labelKey, labelVal)
	flt := "(name=gke-vault*)"
	igms := pkg.FilterInstanceGroupManager(PROJECT_NAME, flt)

	for _, v := range igms {
		fmt.Println(v.Name)
	}
	*/

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
