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
    PKey  bool
}

type TableMetadata map[string]ColumnMetadata

func GetPKey(colData map[string]interface{}) bool {
    pkey, exists := colData["pkey"]
    if !exists {
        return false
    }
    return pkey.(bool)
}

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
            cD := colData.(map[string]interface{})
            newCol := ColumnMetadata{
                Type: cD["type"].(string),
                Name: colName,
                Order: cD["order"].(int),
                PKey: GetPKey(cD),
            }
            cols[colName] = newCol
        }
        tables[tName.(string)] = cols
        fmt.Printf("Read in table %v\n", tName.(string))
    }
    return tables, nil
}
