package di

//go:generate stringer -type=EventTopic -trimprefix=EventTopic -linecomment

type EventTopic int

const (
	EventTopicDIReady EventTopic = iota // di:ready
	EventTopicArgParse // di:arg:parse
)
