package config

import (
    "fmt"
    "io/ioutil"
    
    "gopkg.in/yaml.v3"
)

type ColumnMetadata struct {
    Type  string `json:"type"`
    Name  string 
    Order int
}

type TableMetadata map[string]ColumnMetadata

func ReadConfig(filePath string) (map[string]TableMetadata, error) {
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
