package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "gopkg.in/yaml.v3"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type Server struct {
    Tables map[string]Table
}

func NewServer() *Server {
    return &Server{
        Tables: make(map[string]Table),
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

    fmt.Println(data)

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

    fmt.Println(s.Tables)
}

// TODO: add more metadata
type column struct {
    Type string
}

type Table struct {
    Cols map[string]column
}

func main() {
    s := NewServer()
    s.ReadConfig("tester.yaml")

    e := echo.New()

    e.Use(middleware.Logger())
    e.GET("/accept", s.accept)

    e.Logger.Fatal(e.Start(":1323"))
}

func (s *Server) accept(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}
