package di

//go:generate stringer -type=ServiceRef -output=zz_gen_serviceref_string.go -linecomment

type ServiceRef int

const (
	LoggerService ServiceRef = iota + 1
)
