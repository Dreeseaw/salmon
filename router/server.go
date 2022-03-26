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

func (rs *RoutingServer) SendInsert(ctx context.Context, ic *pb.InsertCommand) (*pb.SuccessResponse, error) {
    return nil, nil
}

func (rs *RoutingServer) SendSelect(sc *pb.SelectCommand, stream pb.RouterService_SendSelectServer) error {
    return nil
}

func (rs *RoutingServer) ReceiveReplicas(stream pb.RouterService_ReceiveReplicasServer) error {
    return nil
}

func (rs *RoutingServer) ProcessPartials(stream pb.RouterService_ProcessPartialsServer) error {
    return nil
}
