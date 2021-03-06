package salmon

import (
    "io"
    "log"
     // "fmt"
//    "time"
    "context"

    "google.golang.org/grpc/metadata"
    
    "github.com/Dreeseaw/salmon/shared/config"
    pb "github.com/Dreeseaw/salmon/shared/grpc"
    cmds "github.com/Dreeseaw/salmon/shared/commands"
)

type ReplicaReceiver struct {
    ClientId    string
    TableData   map[string]config.TableMetadata
    SuccessChan chan cmds.CommandResult
    ManagerChan chan cmds.Command
}

func NewReplicaReceiver(id string, mc chan cmds.Command) *ReplicaReceiver {
    sc := make(chan cmds.CommandResult)
    return &ReplicaReceiver{
        ClientId: id,
        TableData: make(map[string]config.TableMetadata),
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
    go func(mc chan cmds.Command) {
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

            ic := cmds.InsertCommandFromPb(replicaComm, tMeta, rr.SuccessChan)

            // send replica (insert) command
            // fmt.Printf("sending ic\n")
            mc <- &ic
        }
    }(rr.ManagerChan)

    // send success responses back to router
    for {
        select {
        case resp := <- rr.SuccessChan:
            succ := cmds.ResponseToPb(resp)
            if err := stream.Send(succ); err != nil {
                log.Fatalf("Failed to send a note: %v", err)
		    }
        case <-fin:
            return
        }
    }

    return
}
