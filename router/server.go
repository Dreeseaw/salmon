/*
    Main server code & grpc implementation
*/
package main

import (
    "io"
    "fmt"
    "context"
    
    // "google.golang.org/grpc"
    // "github.com/golang/protobuf/proto"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type RoutingServer struct {
    pb.UnimplementedRouterServiceServer

    InsertChan  InsertCommChan // send insert commands to router
    ConnectChan chan *Client // for new clients to be added to engine
}

func NewRoutingServer(ic InsertCommChan, cc chan *Client) *RoutingServer {
    return &RoutingServer{
        InsertChan: ic,
        ConnectChan: cc,
    }
}

// Connect a new client, send over to engine as well
func (rs *RoutingServer) NewClient() *Client{
    // generate new id
    cli := &Client{"test_client", make(InsertCommChan)}
    rs.ConnectChan <- cli
    return cli
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

    // TODO: breakout a unary connect rpc
    // create client (replica channel)
    cli := rs.NewClient()
    fmt.Printf("[server] client connected %v\n", cli.Id)

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
