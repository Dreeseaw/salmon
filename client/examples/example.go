/*
   A very simple usage example

   Requires a running Salmon Router @ localhost:27604
*/
package main

import (
	"fmt"
    "time"
	"os"

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

    // create 4 clients
    var clis [4]*salmon.Salmon
    i := 0
    for i <= 3 {
        sal := salmon.NewSalmon("localhost:27604")
        sal.Init("/tmp/salmon.yaml")
        sal.Start()
        clis[i] = sal
        fmt.Printf("Started salmon client %v\n", i)
        i = i+1
    }

    fmt.Println("gonna insert now")
    i = 0
    for i <= 3 {
        go func(i int) {
            sal := clis[i]
            for j := 0; j < 2; j++ {
                myObj := map[string]interface{}{
                    "testcolumnint": (int32)(i),
                    "testcolumnstr": "tester",
                    "testcolumnfloat": (float64)(j),
                    "testcolumnbool": false,
                }
                err := sal.Insert("testtable", myObj)
                if err != nil {
                    panic(err)
                }
            }
        }(i)
        i = i+ 1
    }
        
    fmt.Println("insert(s) successful, sleeping for 5s")
    time.Sleep(time.Second * 5)

    // print local table
    i = 0
    for i <= 3 {
        sels := []string{"testcolumnint", "testcolumnfloat"}
        objs, err := clis[i].Select("testtable", sels, nil)
        if err != nil {
            panic(err)
        }
        for oi, obj := range objs {
            fmt.Printf("cli=%v,obj=%v,oi=%v\n",i,obj,oi)
        }
        i = i + 1
    }
}
