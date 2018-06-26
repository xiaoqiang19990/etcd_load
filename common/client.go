package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Worker struct {
	Test    string
	Name    string
	IP      string
	Port    string
	KeysAPI client.KeysAPI
}

// workerInfo is the service register information to etcd
type WorkerInfo struct {
	Test string
	Name string
	IP   string
	Port string
}

func NewWorker(test, name, IP, port string, endpoints []string) *Worker {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	w := &Worker{
		Name:    name,
		IP:      IP,
		KeysAPI: client.NewKeysAPI(etcdClient),
		Port:    port,
		Test:    test,
	}
	return w
}
func (w *Worker) GetList() {
	api := w.KeysAPI
	res, err := api.Get(context.Background(), "workers/na1", nil)
	if err != nil {
		fmt.Println("GetList get  error: ", err)
	}
	for _, v := range res.Node.Nodes {
		fmt.Println("GetList get  :", v.Value)
	}
}

func (w *Worker) DelList() {
	api := w.KeysAPI
	_, err := api.Delete(context.Background(), "workers/na1", nil)
	if err != nil {
		fmt.Println("DelList del error: ", err)
	}

}
func (w *Worker) HeartBeat() {
	api := w.KeysAPI
	info := &WorkerInfo{
		Name: w.Name,
		IP:   w.IP,
		Port: w.Port,
	}

	key := "workers/" + w.Name
	value, _ := json.Marshal(info)

	_, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
		TTL: time.Second * 10000,
	})
	if err != nil {
		log.Println("Error update workerInfo:", err)
	}
}
