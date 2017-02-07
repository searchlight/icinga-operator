package util

type IcingaState int32

const (
	Ok       IcingaState = 0
	Warning  IcingaState = 1
	Critical IcingaState = 2
	Unknown  IcingaState = 3
)

var (
	State = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)
