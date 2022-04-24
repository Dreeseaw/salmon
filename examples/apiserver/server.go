/*
    A simple HTTP server that uses Salmon as a datastore

    Primarily designed for whole system testing
*/
package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"

    salmon "github.com/Dreeseaw/salmon/client"
)

type Server struct {
    sal *salmon.Salmon
}

type JsonObject struct {
    Intcol   int     `json:"intcol"`
    Strcol   string  `json:"strcol"`
    Boolcol  bool    `json:"boolcol"`
    Floatcol float64 `json:"floatcol"`
}

func NewServer() *Server {

    sal := salmon.NewSalmon("router:27604")
    sal.Init("config.yaml")

    return &Server{
        sal: sal,
    }
}

func (s *Server) Insert(w http.ResponseWriter, req *http.Request) {
    // get obj from req, insert
    fmt.Println("[apiserver] got insert")
    log.Println("[apiserver] plz?")
    var jo JsonObject
    dec := json.NewDecoder(req.Body)
    err := dec.Decode(&jo)
    if err != nil {
        panic(err) // who cares
    }

    obj := make(map[string]interface{})
    obj["intcol"] = jo.Intcol
    obj["strcol"] = jo.Strcol
    obj["boolcol"] = jo.Boolcol
    obj["floatcol"] = jo.Floatcol
    
    fmt.Println(obj)
    err = s.sal.Insert("maintable", obj)
    if err != nil {
        panic(err)
    }
}

func (s *Server) PrintTable(w http.ResponseWriter, req *http.Request) {
    // dump only this client's table (will have to be changed when selects are distr)
    fmt.Println("[apiserver] got print_table")
    sels := []string{"intcol","strcol","boolcol","floatcol"} // need * identifier
    objs, err := s.sal.Select("maintable", sels, nil)
    if err != nil {
        panic(err) // who cares
    }
    for _, obj := range objs {
        fmt.Println(obj)
    }
}

func main() {
    s := NewServer()

    http.HandleFunc("/insert", s.Insert)
    http.HandleFunc("/print_table", s.PrintTable)

    fmt.Println("Server started on port 8090")
    http.ListenAndServe(":8090", nil)
}
