### check leader status 

`
docker exec etcd-etcd-1-1 etcdctl endpoint status --cluster -w table
`

### put key 

`
docker exec etcd-etcd-1-1 etcdctl put mykey "hello-etcd"
`

### get key

`
docker exec etcd-etcd-3-1 etcdctl get mykey
`