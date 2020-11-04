package config

import (
	"log"
	"time"
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
)

type zkClient struct {
	servers []string
	sessionTimeout time.Duration
}

func NewZKStore() KVStore {
	var hosts = []string{"localhost:2181"}
	var timeout = time.Second*5

	return &zkClient{hosts, timeout}
}

func (c *zkClient) GetConfig(dataId string, group string) (string, error) {
	log.Printf("GetConfig for dataId [%s] group [%s]", dataId, group)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : something terrible happen -> %s ", err)
        return "", err
    }
	defer conn.Close()

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)
	v, _, err := conn.Get(itemPath)
    if err != nil {
        log.Printf("Error : something terrible happen -> %s ", err)
        return "", err
    }

    log.Printf("value of path[%s]=[%s].\n", itemPath, v)
	
	return string(v), nil
}

func (c *zkClient) PublishConfig(dataId string, group string, content string) error {
	log.Printf("PublishConfig for dataId [%s] group [%s] content [%s]", dataId, group, content)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : something terrible happen -> %s", err)
        return err
    }
	defer conn.Close()

	rootPath := fmt.Sprintf("/%s", group)
	exist, _, err := conn.Exists(rootPath)
    if err != nil {
        log.Printf("Error : something terrible happen -> %s", err)
        return err
	}

	var flags int32 = 0
	var acls = zk.WorldACL(zk.PermAll)
	
	if !exist {
		var rootData = []byte("")
		
		p, err_create := conn.Create(rootPath, rootData, flags, acls)
		if err_create != nil {
			log.Printf("Error : something terrible happen -> %s", err_create)
			return err_create
		}
		log.Printf("root node created: %s", p)
	}

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)
	p, err_create := conn.Create(itemPath, []byte(content), flags, acls)
	if err_create != nil {
		log.Printf("Error : something terrible happen -> %s", err_create)
		return err_create
	}
	log.Printf("item node created: %s", p)
	
	return nil
}

func (c *zkClient) MonitorConfig(dataId string, group string, callback ConfigEventCallback) error {
	log.Printf("MonitorConfig for dataId [%s] group [%s]", dataId, group)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : something terrible happen -> %s", err)
        return err
    }
	defer conn.Close()

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)

	for  {
        _, _, ch, _ := conn.GetW(itemPath)
		select {
		case e := <-ch:
			if e.Err == nil {
				if e.Type == zk.EventNodeDataChanged {
                    log.Printf("has node[%s] data changed\n", e.Path)
                    log.Printf("%+v\n", e)
                    v, _, err := conn.Get(itemPath)
                    if err != nil {
                        fmt.Println(err)
                        return
                    }
        
                    fmt.Printf("value of path[%s]=[%s].\n", ch_path, v)
                }
			}
        }
	}
	
	return nil
}