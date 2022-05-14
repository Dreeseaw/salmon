package main

import (
    "fmt"
    "flag"
    "net"

    "google.golang.org/grpc"

    "github.com/Dreeseaw/salmon/shared/config"
    pb "github.com/Dreeseaw/salmon/shared/grpc"
)

type blank struct {}

var (
    port = flag.Int("port", 27604, "the port for the server")
    configPath = flag.String("config", "/etc/salmon.yaml", "path to config file")
)

func main() {
    flag.Parse()

    // get table config for cluster
    tables, err := config.ReadConfig(*configPath)
    if err != nil {
        panic(err)
    }

    // create server/engine channels
    insertChan := make(chan *pb.InsertCommand)

    // universal client map (mutexed)
    clientMap := NewClientMap()

    //start engine
    engine := NewRoutingEngine(
        clientMap,
        insertChan,
        tables,
    )
    go engine.Start()

    //start server
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
}
