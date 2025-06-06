package datatypes

import "time"

type Key struct {
	fld       Type
	schema    *Schema
	readTime  time.Time
	writeTime time.Time
}

func NewKey(fld interface{}, schema *Schema) Key {
	var keyFld Type
	switch fld.(type) {
	case string:
		keyFld = Str{Value: fld.(string)}
	case int:
		keyFld = Integer{Value: fld.(int)}
	}
	return Key{fld: keyFld, schema: schema}
}

func (k *Key) getKeyCopy() Key {
	newK := Key{fld: k.fld, schema: k.schema, readTime: k.readTime, writeTime: k.writeTime}
	return newK
}

func (key *Key) GetField() Type {
	return key.fld
}
