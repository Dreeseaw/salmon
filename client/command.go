package main

import (
    "errors"
//    "strconv"

    pb "github.com/Dreeseaw/salmon/grpc"
)

type Command interface {}

type Object map[string]interface{}

type InsertCommand struct {
    Id         string
    TableName  string
    Obj        Object
    ResultChan chan CommandResult
}

type SelectCommand struct {
    Id         string
    TableName  string
    Selectors  []string
    Filters    []filter
    ResultChan chan CommandResult
}

type CommandResult struct {
    Id      string
    Error   error
    Objects []Object
}

func InsertCommandFromPb(inp *pb.InsertCommand, tm TableMetadata, rc chan CommandResult) *InsertCommand {
    obj, _ := ObjectFromPb(inp.GetObj(), tm)
    return &InsertCommand{
        TableName: inp.GetTable(),
        Obj: obj,
        ResultChan: rc,
    }
}

func ObjectFromPb(inp *pb.Object, tm TableMetadata) (Object, error) {
    obj := make(Object)
    colList := orderColList(tm)

    for i, anyField := range inp.GetField() {
        colName := colList[i].Name
        switch val := anyField.Value.(type) {
        case *pb.FieldType_Sval:
            obj[colName] = val
        case *pb.FieldType_Fval:
            obj[colName] = val
        case *pb.FieldType_Ival:
            obj[colName] = val
        case *pb.FieldType_Bval:
            obj[colName] = val
        case nil:
            return nil, errors.New("nil value found for field")
        default:
            return nil, errors.New("type unknown for field")
        }
    }
    return obj, nil
}

func ResponseToPb(inp CommandResult) *pb.SuccessResponse {
    if inp.Error == nil && inp.Objects == nil {
        return &pb.SuccessResponse{
            Success: true,
            Id: inp.Id,
        }
    }
    return nil
}
