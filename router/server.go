/*
    Main server code & grpc implementation
*/
package main

import (
    "io"
    "fmt"
    // "errors"
    "context"

    "google.golang.org/grpc/metadata"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type idContextKey string

type RoutingServer struct {
    pb.UnimplementedRouterServiceServer

    InsertChan  InsertCommChan // send insert commands to router
    Clients     *ClientMap
}

func NewRoutingServer(cm *ClientMap, ic InsertCommChan) *RoutingServer {
    return &RoutingServer{
        InsertChan: ic,
        Clients: cm,
    }
}

// Unary rpc
func (rs *RoutingServer) Connect(ctx context.Context, ci *pb.ClientID) (*pb.SuccessResponse, error) {
    // TODO: clean up
    cid := ci.GetId()
    cli := &Client{cid, make(InsertCommChan)}
    rs.Clients.Add(cli)
    fmt.Printf("[server] client connected %v\n", cid)
    return &pb.SuccessResponse{Success: true, Id: "connected"}, nil
}

// Unary rpc
func (rs *RoutingServer) Disconnect(ctx context.Context, ci *pb.ClientID) (*pb.SuccessResponse, error) {
    // TODO: disconnect from engine, disburse replicas
    return &pb.SuccessResponse{Success: true, Id: "disconnected"}, nil
}


// Unary rpc
func (rs *RoutingServer) SendInsert(ctx context.Context, ic *pb.InsertCommand) (*pb.SuccessResponse, error) {
    
    // create result chan

    // send to query engine
    fmt.Println("[server] received insert command")
    rs.InsertChan <- ic

    // block for success or not? no

    // send success back
    resp := &pb.SuccessResponse{Success: true, Id: "bob"}
    return resp, nil
}

// Server-side streaming rpc
func (rs *RoutingServer) SendSelect(sc *pb.SelectCommand, stream pb.RouterService_SendSelectServer) error {
    
    // read in select command

    // stream objects back to client

    return nil
}

// Duplex rpc
func (rs *RoutingServer) ReceiveReplicas(stream pb.RouterService_ReceiveReplicasServer) error {
    
    fin := make(chan blank)

    md, _ := metadata.FromIncomingContext(stream.Context())
    cliId := md.Get("id")
    cli, err := rs.Clients.Get(cliId[0])
    if err != nil {
        return err
    }

    // start server streaming goroutine
    go func() {
        for {
            select {
            case <-fin:
                return
            case ic := <-cli.ReplChan:
                fmt.Printf("[server] sending replica to client %v\n", cli.Id)
                if err := stream.Send(ic); err != nil {
                    // TODO: handle error case
                    return
                }
            }
        }
    }()

    // end server streaming goroutine
    defer func() {
        fin <- blank{}
    }()

    for {
        // get success reponses in any order
        succResp, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        if !succResp.GetSuccess() {
            panic("replica failed") // TODO: remove panic
        }
        // TODO: manage failed replica insert better
    }

    return nil
}

// Duplex rpc
func (rs *RoutingServer) ProcessPartials(stream pb.RouterService_ProcessPartialsServer) error {
    return nil
}
