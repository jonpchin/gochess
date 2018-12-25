// +build windows

package gostuff

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func ListServices() {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return
	}

	manager, err := mgr.ConnectRemote(hostname)

	if err != nil {
		fmt.Println(err)
		return
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

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	manager, err := mgr.ConnectRemote(hostname)
	if err != nil {
		fmt.Println("Could not connect to remote host", err)
		return nil
	}
	defer manager.Disconnect()

	mysql, err := manager.OpenService("MySQL80")
	if err != nil {
		fmt.Println("Could not open MySQL80 service", err)
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
			fmt.Println("Could not start MySQL80 service", err)
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
