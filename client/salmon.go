/*
salmon.go

go interface to cache
*/

package main

import (
    "fmt"
    // "errors"
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

type blank struct {}

type ColumnMetadata struct {
    Type  string `json:"type"`
    Name  string 
    Order int
}
type TableMetadata map[string]ColumnMetadata

type Salmon struct {
    ManagerThread  *Manager
    ManagerChannel chan Command
    FinishChannel  chan blank
}

func NewSalmon() *Salmon {
    mc  := make(chan Command)
    fc  := make(chan blank)
    man  := NewManager(ManagerOptions{
        ManChan: mc,
    })

    return &Salmon{
        ManagerThread: man,
        ManagerChannel: mc,
        FinishChannel: fc,
    }
}

// Init the salmon client
func (sal *Salmon) Init() error {

    // read in config file
    tables, err := sal.ReadConfig()
    if err != nil {
        return err
    }

    sal.ManagerThread.Init(tables)
    return nil
}

// Read in a table config yaml
func (sal *Salmon) ReadConfig(filePath string) ([]TableMetadata, error) {
    yfile, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    data := make(map[interface{}]interface{})

    err = yaml.Unmarshal(yfile, &data)
    if err != nil {
        return nil, err
    }

    tables := make(map[string]TableMetadata)

    //TODO: clean up type casting
    for tName, tCols := range data {
        cols := make(TableMetadata)
        for colName, colData := range tCols.(map[string]interface{}) {
            newCol := ColumnMetadata{
                Type: (colData.(map[string]interface{}))["type"].(string),
                Name: colName,
                Order: (colData.(map[string]interface{}))["order"].(int),
            }
            cols[colName] = newCol
        }
        tables[tName.(string)] = cols
        fmt.Println("Read in table %v", tName.(string))
    }
    return tables, nil
}

// Start the salmon client
func (sal *Salmon) Start() error {

    go sal.ManagerThread.Start(sal.FinishChannel)

    return nil
}

// Close client (can be deferred)
func (sal *Salmon) Close() {
    sal.FinishChannel <- blank{}
}

// Insert an object into the system
func (sal *Salmon) Insert(table string, object Object) error {
    
    // create result channel
    rc := make(chan CommandResult)

    // validate object
    cmd := InsertCommand{
        TableName: table,
        Obj: object,
        ResultChan: rc,
    }

    // send to manager channel
    sal.ManagerChannel <- cmd

    // wait for result
    results, _ := <- rc

    return results.Error
}

// Select queries a table in a SQL-ish fashion
func (sal *Salmon) Select(table string, selectors []string, filters []filter) ([]Object, error) {
   
    // result channel
    rc := make(chan CommandResult)

    // create command
    cmd := SelectCommand{
        TableName: table,
        Selectors: selectors,
        Filters: filters,
        ResultChan: rc,
    }

    // send to manchan
    sal.ManagerChannel <- cmd

    // wait for results
    results, _ := <- rc

    if results.Error != nil {
        return results.Objects, nil
    }
    return nil, results.Error

}
