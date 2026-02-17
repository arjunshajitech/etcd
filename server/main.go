// discovery.go
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	services = make(map[string]string)
	mu       sync.RWMutex
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	prefix := "/services/sfu/"

	// 1️⃣ Initial load
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	for _, kv := range resp.Kvs {
		addService(string(kv.Key), string(kv.Value))
	}

	// 2️⃣ Watch prefix
	rch := cli.Watch(ctx, prefix, clientv3.WithPrefix())

	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					addService(string(ev.Kv.Key), string(ev.Kv.Value))
				case clientv3.EventTypeDelete:
					removeService(string(ev.Kv.Key))
				}
			}
		}
	}()

	// Demo: print every 5 seconds
	for {
		time.Sleep(5 * time.Second)
		printServices()
	}
}

func addService(key, val string) {
	mu.Lock()
	defer mu.Unlock()
	services[key] = val
	fmt.Println("Added:", key, val)
}

func removeService(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(services, key)
	fmt.Println("Removed:", key)
}

func printServices() {
	mu.RLock()
	defer mu.RUnlock()

	fmt.Println("Active Services:")
	for k, v := range services {
		fmt.Println("  ", k, "->", v)
	}
}
