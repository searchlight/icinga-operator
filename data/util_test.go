package data_test

import (
	"fmt"
	"log"
	"testing"

	. "github.com/appscode/searchlight/data"
	"github.com/stretchr/testify/assert"
)

func TestLoadClusterChecks(t *testing.T) {
	ic, err := LoadClusterChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ic.Command), "No check agent found")
	fmt.Println(ic.Command[0].Name)
	fmt.Println(ic.Command[0].Description)
}

func TestLoadNodeChecks(t *testing.T) {
	ic, err := LoadNodeChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ic.Command), "No check agent found")
	fmt.Println(ic.Command[0].Name)
	fmt.Println(ic.Command[0].Description)
}

func TestLoadPodChecks(t *testing.T) {
	ic, err := LoadPodChecks()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ic.Command), "No check agent found")
	fmt.Println(ic.Command[0].Name)
	fmt.Println(ic.Command[0].Description)
}
