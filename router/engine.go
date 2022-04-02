/*
    Routing Engine
    The main thread for processing distributed inserts & selects.
    Manages partition placement & creation.
*/

package main

import (

    pb "github.com/Dreeseaw/salmon/grpc"
)

type RoutingEngine struct {
    InsertChan chan *pb.InsertCommand
}


func NewRoutingEngine(ic chan *pb.InsertCommand) *RoutingEngine{
    return &RoutingEngine{
        InsertChan: ic,
    }
}

func (re *RoutingEngine) Start() {

    fin := make(chan blank)

    for {
        select {
            case <-fin:
                return
            case insert := <-re.InsertChan:
                re.ProcessInsert(insert)
        }
    }

    return
}

func (re *RoutingEngine) ProcessInsert(ic *pb.InsertCommand) {
    return
}
