package main

import (
//    "fmt"
    "io/ioutil"
    "net/http"
    "gopkg.in/yaml.v3"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type GUID string
type Partition int
type PartitionList []Partition

type ClientData struct {
    Partitions map[string]PartitionList
}

type Server struct {
    Tables   map[string]Table
    Clients  map[GUID]ClientData
}

func NewServer() *Server {
    return &Server{
        Tables: make(map[string]Table),
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

    for tName, tCols := range data {
        cols := make(map[string]column)
        for colName, colData := range tCols.(map[string]interface{}) {
            newCol := column{
                Type: (colData.(map[string]interface{}))["type"].(string),
            }
            cols[colName] = newCol
        }
        s.Tables[tName.(string)] = Table{Cols: cols}
    }
}

// TODO: add more metadata
type column struct {
    Type string `json:"type"` 
}

type Table struct {
    Cols map[string]column
}

func main() {
    s := NewServer()
    s.ReadConfig("tester.yaml")

    e := echo.New()

    e.Use(middleware.Logger())
    e.GET("/accept", s.Accept)

    e.Logger.Fatal(e.Start(":1323"))
}

// Accept a new client connection, send back table configs 
func (s *Server) Accept(c echo.Context) error {
    // get info from client


    // return table configs
    return c.JSON(http.StatusOK, s.Tables)
}
