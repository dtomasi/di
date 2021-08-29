package di

// StringRef defines a service reference as string.
// As ServiceDef implements fmt.Stringer as an interface for referencing services.
type StringRef string

// String implements fmt.Stringer interface method.
func (r StringRef) String() string {
	return string(r)
}
