package util

const (
	OK       int32 = 0
	WARNING  int32 = 1
	CRITICAL int32 = 2
	UNKNOWN  int32 = 3
)

var (
	State = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)
