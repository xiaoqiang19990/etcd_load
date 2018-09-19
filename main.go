package main

import (
	"etcd_load/common"
	"flag"
	"fmt"
	"time"
)

func main() {
	var role = flag.String("role", "", "master | worker")
	flag.Parse()
	endpoints := []string{"http://127.0.0.1:2379"}
	if *role == "master" {
		/*master := */ common.NewMaster(endpoints)
		//master.WatchWorkers()
	} else if *role == "worker" {
		worker := common.NewWorker("test", "na1", "127.0.0.1", "8080", endpoints)
		worker.HeartBeat()

	} else if *role == "a" {
		worker := common.NewWorker("test", "na2", "127.0.0.2", "8081", endpoints)
		worker.HeartBeat()
	} else if *role == "b" {
		worker := common.NewWorker("test", "na3", "127.0.0.3", "8082", endpoints)
		worker.HeartBeat()
	} else if *role == "c" {
		worker := common.NewWorker("test", "na4", "127.0.0.4", "8083", endpoints)
		worker.HeartBeat()
	} else if *role == "d" {
		worker := common.NewWorker("test", "na5", "127.0.0.5", "8084", endpoints)
		worker.DelList()
	} else if *role == "e" {
		//discovery.GetLists()
	} else {
		fmt.Println("example -h for usage")
	}
	time.Sleep(20 * time.Second)
	fmt.Println("ddddssss")
	select {}
}
