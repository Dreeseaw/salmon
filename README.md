# salmon

An embedded, distrubuted cache built for maximizing data locality.

### Architecture

Designed to store data on the node that'll need it the most, Salmon clients send objects & queries to a centralized Salmon router process (currently single-node), which in turn replicates inserted objects to other clients & plans distributed queries among partitions. 

![salmon architecture](https://github.com/Dreeseaw/salmon/blob/main/salmonarch.png?raw=true)

This makes it very similar to Memcached & Olric (in Embedded Member mode), with the twist of being able to support aggregation queries and route replicated objects to nodes that already own similar data. Salmon attempts to maximize data locality while minimizing replication overhead and aggregation performance. 

### Usage (heavy development)

After cloning the repo, you'll need to define your own tables in a yaml file. Then you can start the routing server.

Right now, only three functions are supported for the client - NewSalmon, Insert, & Select. Select only supports simple filtering and column selecting, with no aggregations supported yet. 


### Developing

#### Starting Router (gRPC server)

~~~
cd router/ && go run ./...
~~~

#### Compiling new gRPC protos

~~~
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/router_service.proto
~~~


