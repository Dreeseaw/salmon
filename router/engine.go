/*
    Routing Engine
    The main thread for processing distributed inserts & selects.
    Manages partition placement & creation.
*/

package main

import (

    pb "github.com/Dreeseaw/salmon/grpc"
)

type Client struct {
    Id string
    ReplChan InsertCommChan
}

type RoutingEngine struct {
    ConnectChan  chan *Client
    InsertInChan InsertCommChan
    ReplRouter   ReplicaRouter
    Clients      map[string]*Client // client id -> Client
}

func NewRoutingEngine(ic InsertCommChan, cc chan *Client) *RoutingEngine{
    return &RoutingEngine{
        ConnectChan: cc,
        InsertInChan: ic,
        ReplRouter: NewBaseReplicaRouter(),
        Clients: make(map[string]*Client),
    }
}

func (re *RoutingEngine) AddClient(cli *Client) {
    re.Clients[cli.Id] = cli
}

func (re *RoutingEngine) Start() {

    fin := make(chan blank)

    for {
        select {
            case <-fin:
                return
            case insert := <-re.InsertInChan:
                re.ProcessInsert(insert)
            case newCli := <-re.ConnectChan:
                re.AddClient(newCli)
        }
    }

    return
}

type InsertCommChan chan *pb.InsertCommand

type ReplChanMap map[string]InsertCommChan

type ReplicaRouter interface {
    GetReplicaList(map[string]*Client) []InsertCommChan
}

type BaseReplicaRouter struct {
}

func NewBaseReplicaRouter() *BaseReplicaRouter {
    return &BaseReplicaRouter{}
}

func (brr *BaseReplicaRouter) GetReplicaList(cliMap map[string]*Client) []InsertCommChan {
    
    // this basic router simply sends the command to all clients
    // for replication, even the sender
    ret := make([]InsertCommChan, 0)

    // simple map -> list of cli repl chans
    for _, cli := range cliMap {
        ret = append(ret, cli.ReplChan)
    }
    return ret
}

func (re *RoutingEngine) ProcessInsert(ic *pb.InsertCommand) {
    
    // get pkeys

    // get list of client replica chans to send to
    replChans := re.ReplRouter.GetReplicaList(re.Clients)

    for _, replChan := range replChans {
        replChan <- ic
    }

    return
}
