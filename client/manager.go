package salmon

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "context"

    "github.com/google/uuid"
    "google.golang.org/grpc"
    // "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    // "github.com/golang/protobuf/proto"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type ManagerOptions struct {
    ManChan    chan Command
    ServerAddr string
}

type Manager struct {
    ClientId     string
    ServerAddr   string
    Tables       map[string]*Table
    ManChan      chan Command
    ReplicaRecv  *ReplicaReceiver
    RouterClient pb.RouterServiceClient
}

func NewManager(mo ManagerOptions) *Manager {
    id := uuid.New()
    cid := strings.Replace(id.String(), "-", "", -1)
    return &Manager{
        ClientId: cid,
        ServerAddr: mo.ServerAddr,
        Tables: make(map[string]*Table),
        ManChan: mo.ManChan,
        ReplicaRecv: NewReplicaReceiver(cid, mo.ManChan),
        RouterClient: nil,
    }
}

// Init manager client & tables
func (m *Manager) Init(tableData map[string]TableMetadata) func() error {

    // init client
    closeFunc, routerCli := m.NewRouterClient()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if _, err := routerCli.Connect(ctx, &pb.ClientID{Id: m.ClientId}); err != nil {
        panic(err)
    }
    m.RouterClient = routerCli

    // init tables
    for tName, tMeta := range tableData {
        table := NewTable(tMeta)
        // TODO: find cleaner way to set up rr tables
        m.ReplicaRecv.TableData[tName] = tMeta
        m.Tables[tName] = table
    }

    return closeFunc
}


// create & start router client
func (m *Manager) NewRouterClient() (func() error, pb.RouterServiceClient) {
   
    if m.ServerAddr == "mock" {
        return func() error {return nil}, NewMockRouterClient() 
    }

    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    conn, err := grpc.Dial(m.ServerAddr, opts...)
    if err != nil {
        panic(err)
    }
    client := pb.NewRouterServiceClient(conn) //type pb.RouterServiceClient
    return conn.Close, client
}


// Start manager
func (m *Manager) Start(fin chan blank) {

    // start receivers
    go m.ReplicaRecv.Start(m.RouterClient)
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
    result := CommandResult{"default", nil, nil}

    // attempt to store locally first (natural object validation)
    table, exists := m.Tables[command.TableName]
    if exists {
        result.Error = table.InsertObject(command.Obj)
    } else {
        result.Error = errors.New(
            fmt.Sprintf("got insert command for unknown table %v", command.TableName),
        )
    }

    if result.Error != nil {
        return result
    }
    fmt.Printf("obj added to %v", m.ClientId)

    // send to router to be replicated
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    resp, err := m.RouterClient.SendInsert(ctx, InsertCommandToPb(command, table.Meta))
    if err != nil {
        result.Error = err
    }

    if !resp.GetSuccess() {
        result.Error = errors.New(
            fmt.Sprintf("Router returned failure for insert id %v, obj stored locally", resp.GetId()),
        )
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
        result.Error = errors.New(
            fmt.Sprintf("got select command for unknown table %v", command.TableName),
        )
    }

    return result
}

/*
    Mocks
*/


type MockRouterClient struct {
}

func NewMockRouterClient() *MockRouterClient {
    return &MockRouterClient{}
}

func (mrc *MockRouterClient) Connect(ctx context.Context, in *pb.ClientID, opts ...grpc.CallOption) (*pb.SuccessResponse, error) {
    return nil, nil
} 

func (mrc *MockRouterClient) Disconnect(ctx context.Context, in *pb.ClientID, opts ...grpc.CallOption) (*pb.SuccessResponse, error) {
    return nil, nil
} 

func (mrc *MockRouterClient) SendInsert(ctx context.Context, in *pb.InsertCommand, opts ...grpc.CallOption) (*pb.SuccessResponse, error) {
    mockResp := &pb.SuccessResponse{Success: true, Id: "default"}
    return mockResp, nil
}

func (mrc *MockRouterClient) SendSelect(ctx context.Context, in *pb.SelectCommand, opts ...grpc.CallOption) (pb.RouterService_SendSelectClient, error) {
    return nil, nil
}

func (mrc *MockRouterClient) ReceiveReplicas(ctx context.Context, opts ...grpc.CallOption) (pb.RouterService_ReceiveReplicasClient, error) {
    mockStream := NewMockStream()
    return mockStream, nil
}

func (mrc *MockRouterClient) ProcessPartials(ctx context.Context, opts ...grpc.CallOption) (pb.RouterService_ProcessPartialsClient, error) {
    return nil, nil
}

type MockStream struct {
    grpc.ClientStream
}

func NewMockStream() *MockStream {
    return &MockStream{}
}

func (ms *MockStream) Send(*pb.SuccessResponse) error {
    return nil
}

func (ms *MockStream) Recv() (*pb.InsertCommand, error) {
    return nil, nil
}

