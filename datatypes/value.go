package datatypes

import "fmt"

type Value struct {
	data   []byte
	schema *Schema
}

const Seperator = "#$"

func NewValue(schema *Schema, vals ...interface{}) (*Value, error) {
	valBytes := make([]byte, 0)
	seperatorBytes := []byte(Seperator)
	for ind, val := range vals {
		fldTypeMap, err := MapTypes(schema.fldTypeNames[ind])
		if err != nil {
			return nil, err
		}
		switch fldTypeMap {
		case StringTypeName:
			fmt.Println(vals[ind])
			strType := Str{Value: val.(string)}
			strBytes, err := strType.IntoBytes()
			if err != nil {
				return nil, err
			}
			valBytes = append(valBytes, strBytes...)

		case IntergerTypeName:
			intType := Integer{Value: val.(int)}
			intBytes, err := intType.IntoBytes()
			if err != nil {
				return nil, err
			}
			valBytes = append(valBytes, intBytes...)
		}
		valBytes = append(valBytes, seperatorBytes...)
	}
	return &Value{data: valBytes, schema: schema}, nil
}

func (v *Value) getValueCopy() Value {
	newV := Value{data: v.data, schema: v.schema}
	return newV
}

func (value *Value) GetData() []byte {
	return value.data
}
