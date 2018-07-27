// +build windows

package gostuff

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func ListServices() {

	manager, err := mgr.ConnectRemote("ux31e")

	if err != nil {
		fmt.Println(err)
	}
	defer manager.Disconnect()
	services, err := manager.ListServices()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, service := range services {
		fmt.Println(service)
	}
}

func getMySQLService() *mgr.Service {

	// Contains host name of Windows machine
	serviceFile := "secret/service.txt"
	host := ReadOneLine(serviceFile)
	if host == "" {
		fmt.Println("Could not find host in", serviceFile, "is the file missing?")
		return nil
	}
	manager, err := mgr.ConnectRemote(host)
	if err != nil {
		fmt.Println("Could not connect to remote host", err)
	}
	defer manager.Disconnect()

	mysql, err := manager.OpenService("MySQL57")
	if err != nil {
		fmt.Println("Could not open MySQL57 service", err)
		return nil
	}
	return mysql
}

func StartMySQLService() {

	mysql := getMySQLService()
	if mysql == nil {
		return
	}

	status, err := mysql.Query()
	if err != nil {
		fmt.Println("Could not query MySQL status", err)
		return
	}

	if status.State == svc.Stopped {
		err = mysql.Start()
		if err != nil {
			fmt.Println("Could not start MySQL57 service", err)
			return
		}
		fmt.Println("MySQL service was started succesfully.")
	}
}

func StopMySQLService() {

	mysql := getMySQLService()
	if mysql == nil {
		return
	}

	status, err := mysql.Control(svc.Stop)
	if err != nil {
		fmt.Println("Can't stop MySQL service.")
		return
	}
	if status.State == svc.Stopped {
		fmt.Println("MySQL service has stopped.")
	} else {
		fmt.Println("MySQL service was not stopped, the status is", status.State)
	}

}
