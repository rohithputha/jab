package datatypes

import (
	"bufio"
	"errors"
	"fmt"
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

func MapTypes(t string) (string, error) {
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
	schema.keyType, err = MapTypes(key[0])
	if err != nil {
		return nil, err
	}
	schema.keyName = key[1]
	fldTypes := make([]string, 0)
	fldNames := make([]string, 0)
	for scanner.Scan() {
		fld := strings.Split(scanner.Text(), ",")
		fldType, err := MapTypes(fld[0])
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

func writeSchemaToFile(filePath string, keyType string, keyName string, fieldTypes []string, fieldNames []string) error {
	if len(fieldTypes) != len(fieldNames) {
		return fmt.Errorf("mismatch in length of field types (%d) and field names (%d)", len(fieldTypes), len(fieldNames))
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create schema file '%s': %w", filePath, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(fmt.Sprintf("%s,%s\n", keyType, keyName))
	if err != nil {
		return fmt.Errorf("failed to write key line: %w", err)
	}
	for i := 0; i < len(fieldTypes); i++ {
		_, err = writer.WriteString(fmt.Sprintf("%s,%s\n", fieldTypes[i], fieldNames[i]))
		if err != nil {
			return fmt.Errorf("failed to write field line for field '%s' (type '%s'): %w", fieldNames[i], fieldTypes[i], err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

func (s *Schema) GetKeyType() string {
	return s.keyType
}

func (s *Schema) GetFld(i int) (fldName string, fldType string) {
	return s.fldNames[i], s.fldTypeNames[i]
}

//---------------------------------------------------------------------

type SchemaMap map[string]*Schema

var SchMap SchemaMap
var once sync.Once

var sm = make(SchemaMap)

func (sm SchemaMap) GetSchema(schemaFile string) (*Schema, error) {
	return sm[schemaFile], nil
}
func (sm SchemaMap) AddSchema(schemaFile string) error {
	s, err := initSchema(schemaFile)
	if err != nil {
		return err
	}
	sm[schemaFile] = s
	return nil
}

func (sm SchemaMap) CreateSchema(schemaFile string, keyName string, keyType string, fldNames []string, fldTypes []string) error {
	return writeSchemaToFile(schemaFile, keyType, keyName, fldTypes, fldNames)
}

func InitSchmaMap() {
	once.Do(func() {
		SchMap = make(SchemaMap)
	})
}
