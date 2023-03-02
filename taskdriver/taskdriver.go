package taskdriver

//
//import (
//	"context"
//	"fmt"
//
//	//"fmt"
//	"strings"
//	"time"
//
//	"github.com/hashicorp/nomad/api"
//	"nomad-cencus/boundary"
//	"nomad-cencus/nomad"
//)
//
//func TaskDriver(nomadAddress string,
//	nomadACLToken string,
//	boundaryAddress string,
//	boundaryOrganization string,
//	boundaryScope string,
//	boundaryAuthMethod string,
//	boundaryCredentials map[string]interface{},
//	hostCatalogName string,
//	hostSetName string,
//	targets []string,
//	pollInterval time.Duration) error {
//
//	// Set up the Nomad and Boundary clients
//	nomadClient, err := nomad.NewClient(nomadAddress, nomadACLToken)
//	if err != nil {
//		return err
//	}
//	boundaryClient, err := boundary.New(boundaryAddress, boundaryOrganization, boundaryScope, boundaryAuthMethod, boundaryCredentials)
//	if err != nil {
//		return err
//	}
//
//	// Get the ID of the host set, creating it if necessary
//
//	//hsc := hostsets.NewClient(boundaryClient.Client)
//	hostCatalog, err := boundaryClient.GetHostCatalogByName(hostCatalogName, boundaryScope)
//	if err != nil {
//		if _, ok := fmt.Errorf(err); ok {
//			hostSet = boundary.NewHostSet()
//			hostSet.Name = hostSetName
//			hostSet.Type = "static"
//			hostSet, err = boundaryClient.HostSets().Create(hostSet)
//			if err != nil {
//				return err
//			}
//		} else {
//			return err
//		}
//	}
//
//	// Keep track of targets that we create so we can delete them if they are no longer running
//	createdTargets := make(map[string]string)
//
//	// Loop indefinitely, updating targets every poll interval
//	for {
//		// Get all running tasks in Nomad
//		tasks, _, err := nomadClient.Jobs().List(context.Background(), &api.QueryOptions{})
//		if err != nil {
//			return err
//		}
//
//		// Update targets for each service in the host set
//		for _, targetName := range targets {
//			targetAddress, err := getTargetAddress(nomadClient, tasks, targetName)
//			if err != nil {
//				return err
//			}
//			if targetAddress != "" {
//				// Create or update the target in Boundary
//				err := targets.CreateOrUpdate(boundaryClient, hostSet.ID, targetName, targetAddress)
//				if err != nil {
//					return err
//				}
//				createdTargets[targetName] = targetAddress
//			} else {
//				// If the task is not running, delete its target from Boundary
//				_, exists := createdTargets[targetName]
//				if exists {
//					err := targets.Delete(boundaryClient, hostSet.ID, targetName)
//					if err != nil {
//						return err
//					}
//					delete(createdTargets, targetName)
//				}
//			}
//		}
//
//		// Wait for the polling interval before checking for updates again
//		time.Sleep(pollInterval)
//	}
//}
//
//func getTargetAddress(nomadClient *api.Client, tasks []*api.TaskListStub, targetName string) (string, error) {
//	// Find the task for the target
//	taskName := getTaskName(targetName)
//	var task *api.TaskListStub
//	for _, t := range tasks {
//		if t.Name == taskName {
//			task = t
//			break
//		}
//	}
//
//	// If the task is running, return its address
//	if task != nil {
//		if len(task.Networks) > 0 && task.Networks[0].IP != "" {
//			return task.Networks[0].IP, nil
//		}
//	}
//
//	// If the task is not running, return an empty string
//	return "", nil
//}
//
//func getTaskName(targetName string) string {
//	parts := strings.Split(targetName, ".")
//	return parts[0]
//}
