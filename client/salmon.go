/*
salmon.go

go interface to cache
*/

package salmon

import (
    // "errors"
    "github.com/Dreeseaw/salmon/shared/config"
    cmds "github.com/Dreeseaw/salmon/shared/commands"
)

type blank struct {}

type Salmon struct {
    ManagerThread  *Manager
    ManagerChannel chan cmds.Command
    FinishChannel  chan blank
    CloseClient    func() error
}

func NewSalmon(serverAddr string) *Salmon {
    mc  := make(chan cmds.Command)
    fc  := make(chan blank)
    man  := NewManager(ManagerOptions{
        ManChan: mc,
        ServerAddr: serverAddr,
    })

    return &Salmon{
        ManagerThread: man,
        ManagerChannel: mc,
        FinishChannel: fc,
        CloseClient: nil,
    }
}

// Init the salmon client
func (sal *Salmon) Init(configFilePath string) error {

    // read in config file
    tables, err := config.ReadConfig(configFilePath)
    if err != nil {
        return err
    }

    sal.CloseClient = sal.ManagerThread.Init(tables)
    return nil
}

// Read in a table config yaml
/*
func (sal *Salmon) ReadConfig(filePath string) (map[string]TableMetadata, error) {
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
        fmt.Printf("Read in table %v\n", tName.(string))
    }
    return tables, nil
}
*/

// Start the salmon client
func (sal *Salmon) Start() error {

    go sal.ManagerThread.Start(sal.FinishChannel)

    return nil
}

// Close salmon client (can be deferred)
func (sal *Salmon) Close() {
    sal.CloseClient()
    sal.FinishChannel <- blank{}
}

// Insert an object into the system
func (sal *Salmon) Insert(table string, object cmds.Object) error {
    
    // create result channel
    rc := make(chan cmds.CommandResult)

    // validate object
    cmd := cmds.InsertCommand{
        Id: "local",
        TableName: table,
        Obj: object,
        ResultChan: rc,
    }

    // send to manager channel
    sal.ManagerChannel <- cmd

    // wait for result
    results, _ := <-rc

    return results.Error
}

// Select queries a table in a SQL-ish fashion
func (sal *Salmon) Select(table string, selectors []string, filters []cmds.Filter) ([]cmds.Object, error) {
   
    // result channel
    rc := make(chan cmds.CommandResult)

    // create command
    cmd := cmds.SelectCommand{
        TableName: table,
        Selectors: selectors,
        Filters: filters,
        ResultChan: rc,
    }

    // send to manchan
    sal.ManagerChannel <- cmd

    // wait for results
    results, _ := <- rc
    
    if results.Error == nil {
        return results.Objects, nil
    }
    return nil, results.Error
}
