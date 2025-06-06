package datatypes

type Type interface {
	IntoBytes() ([]byte, error)
	Equals(Type) bool
	TypeName() string
}
