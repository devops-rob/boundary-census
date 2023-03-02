package main

import (
	"context"
	"fmt"
	"strconv"

	//"github.com/hashicorp/nomad/api"
	nc "nomad-cencus/nomad"
	"os"
)

func main() {
	// Create a new Stream object
	s := nc.NewStream()
	client, _ := nc.NewClient()
	var ipAddr string
	//var port []api.Port

	// Create a context to cancel the event subscription when the program exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to deployment events
	eventStream, err := s.Subscribe(ctx)
	if err != nil {
		s.L.Error("error subscribing to events", "error", err)
		os.Exit(1)
	}

	// Loop over the event stream channel to read events as they happen
	for event := range eventStream {
		for _, e := range event.Events {
			t := e.Type

			switch t {
			case "PlanResult":
				fmt.Println("New Event Received: PlanResult")
			case "JobRegistered":
				fmt.Println("New Event Received: JobRegistered")
				job, _ := e.Job()
				jobId := job.ID
				//metaHcl := job.Meta
				fmt.Println(*jobId)
				allocations, _, _ := client.Jobs().Allocations(*jobId, true, nil)
				for _, alloc := range allocations {
					allocation, _, err := client.Allocations().Info(alloc.ID, nil)
					if err != nil {
						panic(err)
					}

					for _, taskState := range allocation.TaskStates {
						if taskState.State == "running" {
							network := allocation.Resources.Networks
							portsMap := make(map[string]string)
							for _, n := range network {
								ipAddr = n.IP
								dynamicPorts := n.DynamicPorts
								if len(dynamicPorts) > 0 {
									for _, dp := range dynamicPorts {
										strPort := dp.Value
										jobName := allocation.Job.ID
										taskGroupName := allocation.TaskGroup
										targetName := string(*jobName) + `_` + taskGroupName + `_` + string(strPort)
										portsMap[targetName] = strconv.Itoa(strPort)
									}
								}
								reservedPorts := n.ReservedPorts
								if len(reservedPorts) > 0 {
									for _, rp := range reservedPorts {
										strPort := rp.Value
										jobName := allocation.Job.ID
										taskGroupName := allocation.TaskGroup
										targetName := string(*jobName) + `_` + taskGroupName + `_` + strconv.Itoa(strPort)
										portsMap[targetName] = strconv.Itoa(strPort)
									}
								}
							}
							fmt.Printf("IP Address: %v\nPort: %v\n\n", ipAddr, portsMap)
						}
					}
				}
			case "ServiceRegistration":
				service, _ := e.Service()
				serviceID := service.ID
				serviceName := service.ServiceName
				serviceAddr := service.Address
				servicePort := service.Port
				jobID := service.JobID
				allocationID := service.AllocID
				fmt.Printf("service: %s\nID: %s\nip address %s\nport: %s\njob id: %s\nallocation id %s\n", serviceName, serviceID, serviceAddr, strconv.Itoa(servicePort), jobID, allocationID)
			}

		}
	}
}
