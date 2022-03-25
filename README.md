# salmon

An embedded, distrubuted cache built for low-latency, SQL-ish stream processing

WIP: this project was created to play around with distributed computing concepts and should never see the light of day

### Architecture

Designed to store data on the node that'll need it the most, Salmon clients send objects & queries to a centralized Salmon router process (currently single-node), which in turn replicates inserted objects to other clients & plans distributed queries among partitions. This makes it very similar to Memcached & Olric (in Embedded Member mode), with the twist of being able to support aggregation queries and route replicated objects to nodes that already own similar data. Salmon attempts to maximize data locality while minimizing replication overhead and aggregation performance. 

### Usage (heavy development)

After cloning the repo, you'll need to define your own tables in a yaml file. Then you can start the routing server.

Right now, only three functions are supported for the client - NewSalmon, Insert, & Select. Select only supports simple filtering and column selecting, with no aggregations supported yet. 


### Developing

#### Compiling new gRPC protos

~~~
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/router_service.proto
~~~


