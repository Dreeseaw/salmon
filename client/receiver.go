package main

import (
    "io"
    "log"
//    "time"
    "context"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type ReplicaReceiver struct {
    TableData   map[string]TableMetadata
    SuccessChan chan CommandResult
    ManagerChan chan Command
}

func NewReplicaReceiver(mc chan Command) *ReplicaReceiver {
    sc := make(chan CommandResult)
    return &ReplicaReceiver{
        TableData: make(map[string]TableMetadata),
        SuccessChan: sc,
        ManagerChan: mc,
    }
}

func (rr *ReplicaReceiver) Start(client pb.RouterServiceClient) {
    
    ctx := context.Background()
    // defer cancel()

    // create duplex rpc stream
    stream, err := client.ReceiveReplicas(ctx)
    if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}

    fin := make(chan blank)

    // get pb.InsertCommand from router,
    // send InsertCommand to manager
    go func() {
        for {
            replicaComm, err := stream.Recv()
            if err == io.EOF {
                // router closed conn
                fin <- blank{}
                return
            }
            if err != nil {
                log.Fatalf("Failed to receive a replica command : %v", err)
            }

            tMeta, _ := rr.TableData[replicaComm.GetTable()]

            ic := InsertCommandFromPb(replicaComm, tMeta, rr.SuccessChan)

            // send replica (insert) command
            rr.ManagerChan <- ic
        }
    }()

    // send success responses back to router
    for {
        select {
        case resp := <- rr.SuccessChan:
            succ := ResponseToPb(resp)
            if err := stream.Send(succ); err != nil {
                log.Fatalf("Failed to send a note: %v", err)
		    }
        case <-fin:
            return
        }
    }

    return
}
