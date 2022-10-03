package pkg

import (
	"fmt"
	"context"

	"google.golang.org/api/compute/v1"
)

func FilterInstanceGroupManager(projectName, filter string) (rslt []*compute.InstanceGroupManager) {
	rslt = []*compute.InstanceGroupManager{}

	// Get a compute service
	ctx := context.Background()
	cmpsvc, err := compute.NewService(ctx)
	if err != nil {
		fmt.Println("compute NewService error: ", err)
		return
	}

	// Get an InstanceGroupManagersService
	grpMgrSVC := compute.NewInstanceGroupManagersService(cmpsvc)
	instanceGroupManagersAggregatedListCall := grpMgrSVC.AggregatedList(projectName)

	// Do filtering
	instanceGroupManagerAggregatedList, err := instanceGroupManagersAggregatedListCall.Filter(filter).Do()
	if err != nil {
		fmt.Println("Filter InstanceGroupManagerList failed : ", err)
		return
	}


	for _, v := range instanceGroupManagerAggregatedList.Items {
		// igm available under this zone
		if v.Warning == nil {
			rslt = append(rslt, v.InstanceGroupManagers...)
		}
	}

	return
}
