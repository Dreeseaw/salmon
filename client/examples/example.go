/*
   A very simple usage example

   Requires a running Salmon Router @ localhost:27604
*/
package main

import (
	"fmt"
    "time"
	// "os"

	salmon "github.com/Dreeseaw/salmon/client"
)

const (
	configFile string = "/etc/salmon.yaml"
)

func main() {

    // create 4 clients
    var clis [4]*salmon.Salmon

    i := 0
    for i <= 3 {
        sal := salmon.NewSalmon("localhost:27604")
        sal.Init(configFile)
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
            for j := 0; j < 3; j++ {
                myObj := map[string]interface{}{
                    "intcol": (int32)(i),
                    "strcol": "tester",
                    "floatcol": (float64)(j),
                    "boolcol": false,
                }
                err := sal.Insert("maintable", myObj)
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
        sels := []string{"intcol", "floatcol"}
        objs, err := clis[i].Select("maintable", sels, nil)
        if err != nil {
            panic(err)
        }
        for oi, obj := range objs {
            fmt.Printf("cli=%v,obj=%v,oi=%v\n",i,obj,oi)
        }
        i = i + 1
    }
}
