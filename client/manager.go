package main

import (
    "fmt"
    "errors"

    "google.golang.org/grpc"
    // "github.com/golang/protobuf/proto"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type ManagerOptions struct {
    ManChan    chan Command
    ServerAddr string
}

type Manager struct {
    ServerAddr  string
    Tables      map[string]*Table
    ManChan     chan Command
    ReplicaRecv *ReplicaReceiver
}

func NewManager(mo ManagerOptions) *Manager {
    return &Manager{
        ServerAddr: mo.ServerAddr,
        Tables: make(map[string]*Table),
        ManChan: mo.ManChan,
        ReplicaRecv: NewReplicaReceiver(mo.ManChan),
    }
}

// Init manager tables
func (m *Manager) Init(tableData map[string]TableMetadata) {
    for tName, tMeta := range tableData {
        table := NewTable(tMeta)
        m.ReplicaRecv.TableData[tName] = tMeta
        m.Tables[tName] = table
    }
}


// create & start router client
func (m *Manager) NewRouterClient() (func() error, pb.RouterServiceClient) {
   
    if m.ServerAddr == "mock" {
        return func() error {return nil}, NewMockRouterClient() 
    }

    var opts []grpc.DialOption
    conn, err := grpc.Dial(m.ServerAddr, opts...)
    if err != nil {
        panic(err)
    }
    client := pb.NewRouterServiceClient(conn) //type pb.RouterServiceClient
    return conn.Close, client
}


// Start manager
func (m *Manager) Start(fin chan blank) {

    closeFunc, client := m.NewRouterClient()
    defer closeFunc()

    // start receivers
    go m.ReplicaRecv.Start(client)
    // go m.PartialRecv.Start(client)

    // processing loop
    for {
        select {
        case cmd := <- m.ManChan:
            m.Process(cmd)
        case <-fin:
            fmt.Println("Manager thread got shutdown signal")
            return
        }
    }
}

// Process insert or select
func (m *Manager) Process(cmd Command) {
    switch command := cmd.(type) {
    case InsertCommand:
        command.ResultChan <- m.ProcessInsert(command)
    case SelectCommand:
        command.ResultChan <- m.ProcessSelect(command)
    }
}

// Process an Insert command
func (m *Manager) ProcessInsert(command InsertCommand) CommandResult {
    var result CommandResult
    result.Objects = nil

    if table, has := m.Tables[command.TableName]; has {
        result.Error = table.InsertObject(command.Obj)
    } else {
        result.Error = errors.New(fmt.Sprintf("got insert command for unknown table %v", command.TableName))
    }

    return result
}

// Process a select command
func (m *Manager) ProcessSelect(command SelectCommand) CommandResult {
    var result CommandResult

    if table, has := m.Tables[command.TableName]; has {
        resObjects, err := table.Select(command.Selectors, command.Filters)
        if err != nil {
            result.Objects = nil
            result.Error = err
        } else {
            result.Objects = resObjects
            result.Error = nil
        }
    } else {
        result.Objects = nil
        result.Error = errors.New(fmt.Sprintf("got insert command for unknown table %v", command.TableName))
    }

    return result
}
