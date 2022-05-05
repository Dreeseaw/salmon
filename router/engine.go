/*
    Routing Engine
    The main thread for processing distributed inserts & selects.
    Manages partition placement & creation.
*/

package main

import (

    pb "github.com/Dreeseaw/salmon/shared/grpc"
)

type RoutingEngine struct {
    InsertInChan InsertCommChan
    ReplRouter   ReplicaRouter
    Clients      *ClientMap
}

func NewRoutingEngine(cm *ClientMap, ic InsertCommChan) *RoutingEngine{
    return &RoutingEngine{
        InsertInChan: ic,
        ReplRouter: NewBaseReplicaRouter(cm),
        Clients: cm,
    }
}

func (re *RoutingEngine) Start() {

    fin := make(chan blank)

    for {
        select {
            case <-fin:
                return
            case insert := <-re.InsertInChan:
                re.ProcessInsert(insert)
        }
    }

    return
}

type InsertCommChan chan *pb.InsertCommand

type ReplChanMap map[string]InsertCommChan

type ReplicaRouter interface {
    GetReplicaList() []InsertCommChan
}

type BaseReplicaRouter struct {
    ReplicaFactor int
    Clients       *ClientMap
}

func NewBaseReplicaRouter(cm *ClientMap) *BaseReplicaRouter {
    return &BaseReplicaRouter{
        ReplicaFactor: 1,
        Clients: cm,
    }
}

func (brr *BaseReplicaRouter) GetReplicaList() []InsertCommChan {
    
    // this basic router sends the command to all clients
    // for replication, even the sender
    ret := make([]InsertCommChan, 0)

    // simple map -> list of cli repl chans
    //for _, cli := range cliMap {
    //    ret = append(ret, cli.ReplChan)
    //}
    return ret
}

func (re *RoutingEngine) ProcessInsert(ic *pb.InsertCommand) {
    
    // get pkeys

    // get list of client replica chans to send to
    // replChans := re.ReplRouter.GetReplicaList()

    //for _, replChan := range replChans {
    //    replChan <- ic
    //}

    for _, cli := range re.Clients.GetAll() {
        cli.ReplChan <- ic
    }

    return
}
