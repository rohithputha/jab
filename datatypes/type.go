package datatypes

type Type interface {
	ParseValue([]byte)
	IntoBytes() ([]byte, error)
	Equals(Type) bool
	TypeName() string
}
