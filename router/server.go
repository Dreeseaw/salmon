/*
    Main server code & grpc implementation
*/
package main

import (
    "io"
    "fmt"
    "context"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type idContextKey string

type RoutingServer struct {
    pb.UnimplementedRouterServiceServer

    InsertChan  InsertCommChan // send insert commands to router
    ConnectChan chan *Client // for new clients to be added to engine
    ClientMap   map[string]*Client // mirror of engine TODO: extrapolate
}

func NewRoutingServer(ic InsertCommChan, cc chan *Client) *RoutingServer {
    return &RoutingServer{
        InsertChan: ic,
        ConnectChan: cc,
        ClientMap: make(map[string]*Client),
    }
}

// Connect a new client, send over to engine as well
func (rs *RoutingServer) GetClient(id string) *Client {
    if cli, has := rs.ClientMap[id]; has { // TODO: protect map with mutex
        return cli
    }
    return nil 
}

// Unary rpc
func (rs *RoutingServer) Connect(ctx context.Context, ci *pb.ClientID) (*pb.SuccessResponse, error) {
    // TODO: clean up
    cid := ci.GetId()
    cli := &Client{cid, make(InsertCommChan)}
    rs.ConnectChan <- cli
    rs.ClientMap[cid] = cli
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
    cliId := stream.Context().Value(idContextKey("id"))
    cli := rs.GetClient(cliId.(string))

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
