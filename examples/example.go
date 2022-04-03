/*
    A very simple usage example

    Requires a running Salmon Router @ localhost:27604
*/
package main

import (
    "os"
    "fmt"

    salmon "github.com/Dreeseaw/salmon/client"
)

const (
    configFile string = "/tmp/salmon.yaml"
)

func write_config() {

    if _, err := os.Stat(configFile); err == nil {
        os.Remove(configFile)
    }

    yaml_txt := []byte(`testtable:
  testcolumnint:
    type: int
    name: testcolumnint
    order: 0
  testcolumnstr:
    type: string
    name: testcolumnstr
    order: 1
  testcolumnbool:
    type: bool
    name: testcolumnbool
    order: 2
  testcolumnfloat:
    type: float
    name: testcolumnfloat
    order: 3`)
    
    err := os.WriteFile(configFile, yaml_txt, 0644)
    if err != nil {
        panic(err)
    }
}

func main() {
    
    // create table config
    write_config()

    sal := salmon.NewSalmon("localhost:27604")
    sal.Init("/tmp/salmon.yaml")
    sal.Start()
    
    fmt.Println("started salmon client")

}
