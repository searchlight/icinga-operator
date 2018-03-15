package data_test

import (
	"log"
	"testing"

	. "github.com/appscode/searchlight/data"
	"github.com/stretchr/testify/assert"
)

func TestIcingaData(t *testing.T) {
	ic, err := LoadClusterChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ic.Command), "No check agent found")

	in, err := LoadNodeChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(in.Command), "No check agent found")

	ip, err := LoadPodChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ip.Command), "No check agent found")
}
