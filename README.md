# salmon

An embedded, distrubuted cache built for maximizing data locality. Under the hood, the cache uses 

### Architecture

Designed to store data on the node that'll need it the most, Salmon clients send objects & queries to a centralized Salmon router process (currently single-node), which in turn replicates inserted objects to other clients & plans distributed queries among partitions. 

![salmon architecture](./salmonarch.png?raw=true "Salmon's Architecture")

This makes it very similar to Memcached & Olric (in Embedded Member mode), with the twist of being able to support aggregation queries and route replicated objects to nodes that already own similar data. Salmon attempts to maximize data locality while minimizing replication overhead and aggregation performance.

The client-router architecture was chosen over a leader-follower replication scheme, such as using Raft or Bully Election, due to my future goals of exploring complex hashing functions. This complex, custom hashing function may be too much for client applications to performantly process. 

### Usage

Right now, only two functions are supported for the client (other than Init & Start) - Insert and Select. Select only supports simple filtering and column selecting, with no aggregations supported yet. See the example below for usage of these functions.

#### Simple Example

Install Client & Run Router
~~~
$ go get github.com/Dreeseaw/salmon/client 
$ go install github.com/Dreeseaw/salmon/router
$ ./router
Serving on localhost:27604
~~~

Creating a Config
~~~
$ cat /etc/salmon.yaml
sensortable:
  sensorid:
    type: int
    name: sensorid
    order: 0
    pkey: True
  sensortype:
    type: string
    name: sensortype
    order: 1
    pkey: True
  sensorval:
    type: float
    name: sensorval
    order: 2
~~~

Client Application Code
~~~
import (
    "fmt"

    salmon "github.com/Dreeseaw/salmon/client"
)

func main() {
    sal := salmon.NewSalmon("localhost:27604")
    sal.Init("/etc/salmon.yaml")
    sal.Start()

    obj := map[string]interface{
        "sensorid": 21,
        "sensortype": "street",
        "sensorval": 12.345,
    }
    sal.Insert("sensortable", obj)

    filters := []Filter{
        StringFilter{
            Col: "sensortype",
            Op: "=",
            Val: "street",
        }
    }
    selectors := []string{"sensorid", "sensorval"}
    resultsObjects, err := sal.Select("sensortable", selectors, filters)

    fmt.Println(resultObjects)
}

~~~

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


