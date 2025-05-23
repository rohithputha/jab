package datatypes

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
)

type Schema struct {
	keyType        string
	keyName        string
	fldTypeNames   []string
	fldNames       []string
	additionalKeys []int
	totalFlds      int
}

type SchemaFace interface {
	//InitSchema(schemaFile string) (*Schema, error)
	GetKeyType() string
	GetFld(fldNumber int) (fldName string, fldType string)
}

func mapTypes(t string) (string, error) {
	switch t {
	case StringTypeName:
		return StringTypeName, nil
	case IntergerTypeName:
		return IntergerTypeName, nil
	}
	return "", errors.New("type not supported")
}

func initSchema(schemaFile string) (*Schema, error) {
	file, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	schema := Schema{}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fLine := scanner.Text()
	key := strings.Split(fLine, ",")
	schema.keyType, err = mapTypes(key[0])
	if err != nil {
		return nil, err
	}
	schema.keyName = key[1]
	fldTypes := make([]string, 0)
	fldNames := make([]string, 0)
	for scanner.Scan() {
		fld := strings.Split(scanner.Text(), ",")
		fldType, err := mapTypes(fld[0])
		fldName := fld[1]
		if err != nil {
			return nil, err
		}
		fldTypes = append(fldTypes, fldType)
		fldNames = append(fldNames, fldName)
	}
	schema.fldTypeNames = fldTypes
	schema.fldNames = fldNames
	schema.totalFlds = len(fldTypes)
	return &schema, nil
}

func (s *Schema) GetKeyType() string {
	return s.keyType
}

func (s *Schema) GetFld(i int) (fldName string, fldType string) {
	return s.fldNames[i], s.fldTypeNames[i]
}

//---------------------------------------------------------------------

type SchemaMap map[string]*Schema

var SchMap map[string]*Schema
var once sync.Once

var sm = make(SchemaMap)

func (sm SchemaMap) GetSchema(schemaFile string) (*Schema, error) {
	return sm[schemaFile], nil
}
func (sm SchemaMap) addSchema(schemaFile string) error {
	s, err := initSchema(schemaFile)
	if err != nil {
		return err
	}
	sm[schemaFile] = s
	return nil
}

func InitSchmaMap() {
	once.Do(func() {
		SchMap = make(SchemaMap)
	})
}
