package main

//import (
//	"time"
//
//	"github.com/devops-rob/taskdriver"
//)
//
//func main() {
//	// Set up the Nomad and Boundary clients
//	nomadAddress := "http://localhost:4646"
//	nomadACLToken := ""
//	boundaryAddress := "http://localhost:9200"
//	boundaryOrganization := "example"
//	boundaryScope := "example"
//	boundaryAuthMethod := "local"
//	boundaryCredentials := map[string]interface{}{
//		"username": "admin",
//		"password": "password",
//	}
//
//	// Set up the host set and targets to create in Boundary
//	hostSetName := "example-host-set"
//	targets := []string{"web", "db"}
//
//	// Set up the polling interval for updating targets
//	pollInterval := 10 * time.Second
//
//	// Start the task driver
//	err := taskdriver.TaskDriver(nomadAddress, nomadACLToken, boundaryAddress, boundaryOrganization, boundaryScope, boundaryAuthMethod, boundaryCredentials, hostSetName, targets, pollInterval)
//	if err != nil {
//		panic(err)
//	}
//}
