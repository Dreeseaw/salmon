package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "gopkg.in/yaml.v3"
    "net"


    "google.golang.org/grpc"

    pb "github.com/Dreeseaw/salmon/shared/grpc"
)

type blank struct {}

func ReadConfig(filePath string) map[string]TableMetadata {
    yfile, err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    data := make(map[interface{}]interface{})

    err = yaml.Unmarshal(yfile, &data)
    if err != nil {
        panic(err)
    }

    tables := make(map[string]TableMetadata)

    //TODO: clean up type casting
    for tName, tCols := range data {
        cols := make(TableMetadata)
        for colName, colData := range tCols.(map[string]interface{}) {
            newCol := ColumnMetadata{
                Type: (colData.(map[string]interface{}))["type"].(string),
            }
            cols[colName] = newCol
        }
        tables[tName.(string)] = cols
    }
    return tables
}

// TODO: add more metadata
type ColumnMetadata struct {
    Type string `json:"type"` 
}

type TableMetadata map[string]ColumnMetadata

var (
    port = flag.Int("port", 27604, "the port for the server")
)

func main() {
    flag.Parse()

    // create server/engine channels
    insertChan := make(chan *pb.InsertCommand)

    // universal client map (mutexed)
    clientMap := NewClientMap()

    //start engine
    engine := NewRoutingEngine(clientMap, insertChan)
    go engine.Start()

    //start server
    //go func(){
    var opts []grpc.ServerOption
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        panic(err)
    }
    grpcServer := grpc.NewServer(opts...)
    pb.RegisterRouterServiceServer(
        grpcServer, 
        NewRoutingServer(clientMap, insertChan),
    )
    fmt.Printf("Serving on localhost:%d\n", *port)
    grpcServer.Serve(lis)
    //}()
}
