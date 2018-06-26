package common

import (
	"encoding/json"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

type Master struct {
	members map[string][]Member
	KeysAPI client.KeysAPI
}

// Member is a client machine
type Member struct {
	InGroup bool
	IP      string
	Name    string
	Port    string
	Test    string
}

var Memberss map[string][]Member

func NewMaster(endpoints []string) *Master {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	master := &Master{
		members: make(map[string][]Member),
		KeysAPI: client.NewKeysAPI(etcdClient),
	}

	go master.WatchWorkers()
	return master
}

var mem []Member

func (m *Master) AddWorker(info *WorkerInfo) {
	member := Member{
		InGroup: true,
		IP:      info.IP,
		Name:    info.Name,
		Port:    info.Port,
		Test:    info.Test,
	}
	mem = append(mem, member)
	m.members[member.Test] = mem

	if _, ok := Memberss[info.Test]; ok {
		Memberss[info.Test] = mem
	} else {
		Memberss = m.members
	}

}
func (m *Master) DelWorker(info *WorkerInfo) {
	del := make([]Member, 0)
	if v, ok := m.members[info.Test]; ok {
		for _, vs := range v {
			if vs.IP == info.IP {
				continue
			}
			del = append(del, vs)
		}
	}
	mem = del
	m.members[info.Test] = mem

	Memberss[info.Test] = m.members[info.Test]

}
func (m *Master) UpdateWorker(info *WorkerInfo) {
	_, _ = m.members[info.Name]
	// if ok {
	// 	member.InGroup = true
	// }

}

func NodeToWorkerInfo(node *client.Node, str string) *WorkerInfo {
	log.Println(node.Value, str)
	info := &WorkerInfo{}
	err := json.Unmarshal([]byte(node.Value), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (m *Master) WatchWorkers() {
	api := m.KeysAPI
	watcher := api.Watcher("workers/", &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers:", err)
			break
		}

		switch res.Action {
		case "set":
			info := NodeToWorkerInfo(res.Node, "set ")
			if member, ok := m.members[info.Name]; ok {
				for _, v := range member {
					if v.IP == info.IP {
						//	log.Println("Update worker ", info.IP)
						m.UpdateWorker(info)
					} else {
						//	log.Println("Add1 worker ", info.IP)
						m.AddWorker(info)
					}
				}

			} else {
				//log.Println("Add2 worker ", info.IP)
				m.AddWorker(info)
			}
		case "delete", "expire":
			info := NodeToWorkerInfo(res.PrevNode, "del")
			//	fmt.Println("Delete worker---------", info.Name)
			m.DelWorker(info)

		}
	}
}
