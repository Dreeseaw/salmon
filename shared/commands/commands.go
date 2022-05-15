package commands

import (
    "errors"

    "github.com/Dreeseaw/salmon/shared/config"
    pb "github.com/Dreeseaw/salmon/shared/grpc"
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

func InsertCommandFromPb(inp *pb.InsertCommand, tm config.TableMetadata, rc chan CommandResult) InsertCommand {
    obj, err := ObjectFromPb(inp.GetObj(), tm)
    if err != nil {
        panic(err)
    }
    return InsertCommand{
        Id: "fromrouter",
        TableName: inp.GetTable(),
        Obj: obj,
        ResultChan: rc,
    }
}

func InsertCommandToPb(inp InsertCommand, tm config.TableMetadata) *pb.InsertCommand {

    fields := make([]*pb.FieldType, len(tm))

    // TODO: factor out a ObjectToPb func
    for colName, colMeta := range tm {
        val, _ := inp.Obj[colName]
        field := new(pb.FieldType)
        
        switch colMeta.Type {
        case "string":
            field.Value = &pb.FieldType_Sval{Sval: val.(string)}
        case "int":
            field.Value = &pb.FieldType_Ival{Ival: val.(int32)}
        case "bool":
            field.Value = &pb.FieldType_Bval{Bval: val.(bool)}
        case "float":
            field.Value = &pb.FieldType_Fval{Fval: val.(float64)}
        default:
            panic("honestly how?")
        }

        fields[colMeta.Order] = field
    }

    return &pb.InsertCommand{
        Iid: inp.Id,
        Table: inp.TableName,
        Obj: &pb.Object{Field: fields},
    }
}

func ObjectFromPb(inp *pb.Object, tm config.TableMetadata) (Object, error) {
    obj := make(Object)
    colList := config.OrderColList(tm)

    for i, anyField := range inp.GetField() {
        colName := colList[i].Name
        switch val := anyField.Value.(type) {
        case *pb.FieldType_Sval:
            obj[colName] = val.Sval
        case *pb.FieldType_Fval:
            obj[colName] = val.Fval
        case *pb.FieldType_Ival:
            obj[colName] = val.Ival
        case *pb.FieldType_Bval:
            obj[colName] = val.Bval
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
