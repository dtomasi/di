package greeter

//go:generate stringer -type=ServiceRef

type ServiceRef int

const (
	ServiceGreeterMorning ServiceRef = iota
	ServiceGreeterAfternoon
	ServiceGreeterEvening
)
