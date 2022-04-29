package di

//go:generate stringer -type=EventTopic -trimprefix=EventTopic -output=zz_gen_eventtopic_string.go -linecomment

type EventTopic int

const (
	EventTopicDIReady EventTopic = iota // di:ready
)
