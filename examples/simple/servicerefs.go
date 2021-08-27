package main

//go:generate stringer -type=ServiceRef

type ServiceRef int

const (
	ServiceGoodMorningGreeter ServiceRef = iota
	ServiceGoodAfternoonGreeter
	ServiceGoodEveningGreeter
)
