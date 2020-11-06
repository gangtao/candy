### preparation
install go https://golang.org/doc/install

install protobuf https://grpc.io/docs/protoc-installation/

### issue
due to the dependency break of grpc/protobuf/etcd-clientv3, refer to this https://github.com/etcd-io/etcd/pull/11564 we cannot use both the latest protobuf and etcd client, need wait for clientv3@v3.5