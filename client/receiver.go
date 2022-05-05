package salmon

import (
    "io"
    "log"
     // "fmt"
//    "time"
    "context"

    "google.golang.org/grpc/metadata"
    
    pb "github.com/Dreeseaw/salmon/shared/grpc"
)

type ReplicaReceiver struct {
    ClientId    string
    TableData   map[string]TableMetadata
    SuccessChan chan CommandResult
    ManagerChan chan Command
}

func NewReplicaReceiver(id string, mc chan Command) *ReplicaReceiver {
    sc := make(chan CommandResult)
    return &ReplicaReceiver{
        ClientId: id,
        TableData: make(map[string]TableMetadata),
        SuccessChan: sc,
        ManagerChan: mc,
    }
}

func (rr *ReplicaReceiver) Start(client pb.RouterServiceClient) {
    
    md := metadata.New(map[string]string{"id": rr.ClientId})
    ctx := metadata.NewOutgoingContext(context.Background(), md)
    // defer cancel()

    // create duplex rpc stream
    stream, err := client.ReceiveReplicas(ctx)
    if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}

    fin := make(chan blank)

    // get pb.InsertCommand from router,
    // send InsertCommand to manager
    go func(mc chan Command) {
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
            // fmt.Printf("sending ic\n")
            mc <- &ic
        }
    }(rr.ManagerChan)

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
