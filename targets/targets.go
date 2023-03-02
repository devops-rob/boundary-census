package targets

//
//import (
//	"github.com/hashicorp/boundary/api/boundary"
//)
//
//func CreateOrUpdate(boundaryClient *boundary.Client, hostSetID string, targetName string, targetAddress string) error {
//	// Get the current target in the host set, if it exists
//	hostSet, err := boundaryClient.HostSets().Get(hostSetID)
//	if err != nil {
//		return err
//	}
//	var target *boundary.Target
//	for _, t := range hostSet.Targets {
//		if t.Name == targetName {
//			target = t
//			break
//		}
//	}
//
//	// If the target does not exist, create it
//	if target == nil {
//		newTarget := boundary.NewTarget()
//		newTarget.Name = targetName
//		newTarget.Type = "tcp"
//		newTarget.Address = targetAddress
//
//		target, err = boundaryClient.Targets().Create(newTarget)
//		if err != nil {
//			return err
//		}
//	}
//
//	// Update the target if its address has changed
//	if target.Address != targetAddress {
//		target.Address = targetAddress
//		_, err = boundaryClient.Targets().Update(target)
//		if err != nil {
//			return err
//		}
//	}
//
//	// Add the target to the host set if it is not already there
//	if !contains(hostSet.TargetIDs, target.ID) {
//		hostSet.TargetIDs = append(hostSet.TargetIDs, target.ID)
//		err = boundaryClient.HostSets().Update(hostSet)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func Delete(boundaryClient *boundary.Client, hostSetID string, targetName string) error {
//	// Get the target from the host set, if it exists
//	hostSet, err := boundaryClient.HostSets().Get(hostSetID)
//	if err != nil {
//		return err
//	}
//	var target *boundary.Target
//	for _, t := range hostSet.Targets {
//		if t.Name == targetName {
//			target = t
//			break
//		}
//	}
//
//	// If the target exists, delete it from the host set and Boundary
//	if target != nil {
//		hostSet.TargetIDs = remove(hostSet.TargetIDs, target.ID)
//		err = boundaryClient.HostSets().Update(hostSet)
//		if err != nil {
//			return err
//		}
//		err = boundaryClient.Targets().Delete(target.ID)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func contains(arr []string, str string) bool {
//	for _, a := range arr {
//		if a == str {
//			return true
//		}
//	}
//	return false
//}
//
//func remove(arr []string, str string) []string {
//	index := -1
//	for i, s := range arr {
//		if s == str {
//			index = i
//			break
//		}
//	}
//	if index >= 0 {
//		arr = append(arr[:index], arr[index+1:]...)
//	}
//	return arr
//}
