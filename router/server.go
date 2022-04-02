/*
    Main server code & grpc implementation
*/
package main

import (
    "io"
    "context"
    
    // "google.golang.org/grpc"
    // "github.com/golang/protobuf/proto"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type RoutingServer struct {
    pb.UnimplementedRouterServiceServer

    InsertChan chan *pb.InsertCommand
    ClientReplicaChan chan *pb.InsertCommand
}

func NewRoutingServer(ic chan *pb.InsertCommand) *RoutingServer {
    s := &RoutingServer{
        InsertChan: ic,
    }
    return s
}

// Unary rpc
func (rs *RoutingServer) SendInsert(ctx context.Context, ic *pb.InsertCommand) (*pb.SuccessResponse, error) {
    
    // create result chan

    // send to query engine

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

    // create replica channel for requesting client
    

    // start server streaming goroutine
    go func() {
        for {
            select {
            case <-fin:
                return
            case ic := <-rs.ClientReplicaChan:
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
        _, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        
        // send to routing engine

    }

    return nil
}

// Duplex rpc
func (rs *RoutingServer) ProcessPartials(stream pb.RouterService_ProcessPartialsServer) error {
    return nil
}
