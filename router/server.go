/*
    Main server code & grpc implementation
*/
package main

import (
    "context"
    
    // "google.golang.org/grpc"
    // "github.com/golang/protobuf/proto"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type RoutingServer struct {
    pb.UnimplementedRouterServiceServer

     
}

func NewRoutingServer() *RoutingServer {
    s := &RoutingServer{}
    return s
}

// Unary rpc
func (rs *RoutingServer) SendInsert(ctx context.Context, ic *pb.InsertCommand) (*pb.SuccessResponse, error) {
    
    // create result chan

    // send to query engine

    // block for success or not?

    // send success back
    resp := &pb.SuccessResponse{nil, nil}
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
    
    fin := chan blank

    // create replica channel for requesting client
    

    // start server streaming goroutine
    go func() {
        for {
            select {
            case <-fin:
                return
            case ic <-ClientReplicaChan:
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
        in, err := stream.Recv()
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
