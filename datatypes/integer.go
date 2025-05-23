package datatypes

import (
	"bytes"
	"encoding/binary"
)

const IntergerTypeName = "integer"

type Integer struct {
	Value int
}

func (i *Integer) ParseValue(bytes []byte) {
	i.Value = int(binary.LittleEndian.Uint64(bytes))
}

func (i *Integer) IntoBytes() ([]byte, error) {
	b := make([]byte, 8)
	buf := bytes.NewBuffer(b)
	err := binary.Write(buf, binary.LittleEndian, i.Value)
	return b, err
}

func (i *Integer) Equals(t Type) bool {
	o, ok := t.(*Integer)
	return ok && i.Value == o.Value
}

func (i *Integer) TypeName() string {
	return IntergerTypeName
}
