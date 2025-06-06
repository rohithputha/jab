package memtbl

import (
	"github.com/rohithputha/jab/datatypes"
	"sync"
	"time"
	_ "time"
)

type Tbl struct {
	kv      map[datatypes.Key]*datatypes.Value
	schema  *datatypes.Schema
	walChan chan<- []byte
}

type Tables struct {
	NameTblMap map[string]*Tbl
	walChan    chan<- []byte
	wal        *wal
}

var TablesObj Tables
var once sync.Once

type TableFace interface {
	AddKV(*datatypes.Key, *datatypes.Value) error
	CreateAndAddKV(interface{}, ...interface{}) error
	GetKV(*datatypes.Key) *datatypes.Value
	CreateAndGetKV(interface{}) *datatypes.Value
}

func CreateTable(name string, keyName string, keyType string, fldNames []string, fldTypes []string) error {
	datatypes.InitSchmaMap()
	schemaName := name + ".sch"
	err := datatypes.SchMap.CreateSchema(schemaName, keyName, keyType, fldNames, fldTypes)
	if err != nil {
		return err
	}

	err = datatypes.SchMap.AddSchema(schemaName)
	if err != nil {
		return err
	}
	schema, err := datatypes.SchMap.GetSchema(schemaName)
	if err != nil {
		return err
	}
	tbl := Tbl{
		kv:      make(map[datatypes.Key]*datatypes.Value),
		schema:  schema,
		walChan: TablesObj.walChan,
	}
	TablesObj.NameTblMap[name] = &tbl
	return nil
}

func getLogString(key *datatypes.Key, val *datatypes.Value) ([]byte, error) {
	currenTime := time.Now()
	timestampVal := datatypes.Integer64{Value: currenTime.Unix()}
	log, err := timestampVal.IntoBytes64()
	if err != nil {
		return nil, err
	}
	keyByes, err := key.GetField().IntoBytes()
	log = append(log, byte('#'))
	log = append(log, byte('$'))
	log = append(log, keyByes...)
	if err != nil {
		return nil, err
	}
	log = append(log, byte('#'))
	log = append(log, byte('$'))
	log = append(log, val.GetData()...)
	return log, nil
}

func (t *Tbl) AddKV(key *datatypes.Key, val *datatypes.Value) error {
	log, err := getLogString(key, val)
	if err != nil {
		return err
	}
	t.walChan <- log
	t.kv[(*key)] = val
	return nil
}

func (t *Tbl) CreateAndAddKV(key interface{}, vals ...interface{}) error {
	keyType := datatypes.NewKey(key, t.schema)
	valType, err := datatypes.NewValue(t.schema, vals...)
	if err != nil {
		return err
	}
	return t.AddKV(&keyType, valType)
}
func (t *Tbl) GetKV(key *datatypes.Key) (val *datatypes.Value) {
	return t.kv[*key]
}

func (t *Tbl) CreateAndGetKV(key interface{}) *datatypes.Value {
	keyType := datatypes.NewKey(key, t.schema)
	return t.GetKV(&keyType)
}

func InitTables(latestLogFileName string) {
	walChan := make(chan []byte)
	wal := InitWal(walChan, latestLogFileName)
	once.Do(func() {
		TablesObj = Tables{
			NameTblMap: make(map[string]*Tbl),
			walChan:    walChan,
			wal:        wal,
		}
		go wal.ListenToLogs()
	})
}
