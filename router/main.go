package main

import (
//    "fmt"
    "flag"
    "io/ioutil"
    "net/http"
    "gopkg.in/yaml.v3"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type blank struct {}
type GUID string
type Partition int
type PartitionList []Partition

type ClientData struct {
    Partitions map[string]PartitionList
}

type Server struct {
    Tables   map[string]TableMetadata
    Clients  map[GUID]ClientData
}

func NewServer() *Server {
    return &Server{
        Tables: make(map[string]TableMetadata),
        Clients: make(map[GUID]ClientData),
    }
}

func (s *Server) ReadConfig(filePath string) {
    yfile, err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    data := make(map[interface{}]interface{})

    err = yaml.Unmarshal(yfile, &data)
    if err != nil {
        panic(err)
    }

    //TODO: clean up type casting
    for tName, tCols := range data {
        cols := make(TableMetadata)
        for colName, colData := range tCols.(map[string]interface{}) {
            newCol := ColumnMetadata{
                Type: (colData.(map[string]interface{}))["type"].(string),
            }
            cols[colName] = newCol
        }
        s.Tables[tName.(string)] = cols
    }
}

// TODO: add more metadata
type ColumnMetadata struct {
    Type string `json:"type"` 
}

type TableMetadata map[string]ColumnMetadata


func main() {
    s := NewServer()
    s.ReadConfig("tester.yaml")

    e := echo.New()

    e.Use(middleware.Logger())
    e.GET("/accept", s.Accept)

    e.Logger.Fatal(e.Start(":1323"))
}

var (
    port = flag.Int("port", 27604, "the port for the server")
)

func grpc_main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)
    pb.RegisterRoutingServiceServer(grpcServer, newRoutingServer())
    grpcServer.Serve(lis)
}

// Accept a new client connection, send back table configs 
func (s *Server) Accept(c echo.Context) error {
    // get info from client


    // return table configs
    return c.JSON(http.StatusOK, s.Tables)
}
