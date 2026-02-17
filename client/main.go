// register.go

// register_session.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("provide node name")
	}
	nodeID := os.Args[1]
	key := "/services/sfu/" + nodeID
	value := "127.0.0.1:9000"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Create session (TTL 10 seconds)
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Println("Session lease:", session.Lease())

	// Register service with session lease
	_, err = cli.Put(
		context.Background(),
		key,
		value,
		clientv3.WithLease(session.Lease()),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registered:", key)

	// Wait until session expires
	<-session.Done()
	fmt.Println("Session expired")
}


// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	clientv3 "go.etcd.io/etcd/client/v3"
// )

// func main() {
// 	nodeID := os.Args[1] // pass node name
// 	key := "/services/sfu/" + nodeID
// 	value := "127.0.0.1:9000"

// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   []string{"localhost:2379"},
// 		DialTimeout: 5 * time.Second,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cli.Close()

// 	ctx := context.Background()

// 	// Create lease
// 	leaseResp, err := cli.Grant(ctx, 10)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Register service
// 	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Registered:", key)

// 	// Keep alive
// 	ch, err := cli.KeepAlive(ctx, leaseResp.ID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for ka := range ch {
// 		fmt.Println("Lease TTL:", ka.TTL)
// 	}
// }
