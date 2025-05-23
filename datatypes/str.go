package datatypes

const StringTypeName = "string"

type Str struct {
	Value string
}

func (s *Str) ParseValue(bytes []byte) {
	s.Value = string(bytes)
}

func (s *Str) IntoBytes() ([]byte, error) {
	return []byte(s.Value), nil
}

func (s *Str) Equals(t Type) bool {
	o, ok := t.(*Str)
	return ok && s.Value == o.Value
}

func (s *Str) TypeName() string {
	return StringTypeName
}

//func (s *Str) setValue(value string) {
//	s.value = value
//}
