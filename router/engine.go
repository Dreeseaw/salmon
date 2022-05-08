/*
    Routing Engine
    The main thread for processing distributed inserts & selects.
    Manages partition placement & creation.
*/

package main

import (
    "errors"

    "github.com/Dreeseaw/salmon/shared/config"
    pb "github.com/Dreeseaw/salmon/shared/grpc"
)

type InsertCommChan chan *pb.InsertCommand

type RoutingEngine struct {
    InsertInChan InsertCommChan
    Clients      *ClientMap
    Tables       map[string]config.TableMetadata
    Partitions   map[string]*PartitionSet
}

func NewRoutingEngine(cm *ClientMap, ic InsertCommChan, tm map[string]config.TableMetadata) *RoutingEngine {
    ps := make(map[string]*PartitionSet)
    ids := make([]string, 0)

    for _, cli := range cm.GetAll() {
        ids = append(ids, cli.Id)
    }
    
    for tName, tMeta := range tm {
        pks := make([]PKey, 0)
        for cName, cMeta := range tMeta {
            if cMeta.PKey {
                pks = append(pks, NewPKey(cName, cMeta.Type))
            }
        }
        ps[tName] = NewPartitionSet(pks, ids)
    }

    return &RoutingEngine{
        InsertInChan: ic,
        Clients: cm,
        Tables: tm,
        Partitions: ps,
    }
}

func (re *RoutingEngine) Start() {
    fin := make(chan blank)

    for {
        select {
            case <-fin:
                return
            case insert := <-re.InsertInChan:
                re.ProcessInsert(insert)
        }
    }

    return
}

type ReplChanMap map[string]InsertCommChan

func (re *RoutingEngine) ProcessInsert(ic *pb.InsertCommand) {
    
    // get table config for plucking pkeys
    tName := ic.GetTable()
    tMeta, exists := re.Tables[tName]
    if !exists {
        panic("got insert for inexistent table") // TODO: return error, not panic
    }
    originClientID := ic.GetIid()

    // find partition set (existence gaurenteed)
    pSet, _ := re.Partitions[ic.GetTable()]

    pSet.UpdateClients(re.Clients)

    obj, err := ObjectFromPb(ic.GetObj(), tMeta)
    if err != nil {
        panic(err) // TODO: return error, not panic
    }
   
    replCliIds := pSet.Process(obj, originClientID)

    clis, _ := re.Clients.GetMany(replCliIds)
    for _, cli := range clis { 
        if originClientID != cli.Id { // origin already owns partition/object
            cli.ReplChan <- ic
        }
    }

    return
}

func orderColList(tm config.TableMetadata) []config.ColumnMetadata {
    ret := make([]config.ColumnMetadata, len(tm))
    for _, colMeta := range tm {
        ret[colMeta.Order] = colMeta
    }
    return ret
} 

// TODO: create either utils or commands shared pkg
// ^ would require refactoring some types in client
func ObjectFromPb(inp *pb.Object, tm config.TableMetadata) (map[string]interface{}, error) {
    obj := make(map[string]interface{})
    colList := orderColList(tm)

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
