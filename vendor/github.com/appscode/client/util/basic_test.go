package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicEncode(t *testing.T) {
	encoded := BasicEncode("appscode", "sadlil", "testpassword")
	assert.Equal(t, "YXBwc2NvZGUuc2FkbGlsOnRlc3RwYXNzd29yZA==", encoded)

	encoded = BasicEncode("@pps^ode", "sadlil", "testpassword")
	assert.Equal(t, "QHBwc15vZGUuc2FkbGlsOnRlc3RwYXNzd29yZA==", encoded)
}
