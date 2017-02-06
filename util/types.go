package util

type IcingaState int32

const (
	OK       IcingaState = 0
	WARNING  IcingaState = 1
	CRITICAL IcingaState = 2
	UNKNOWN  IcingaState = 3
)

var (
	State = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)
