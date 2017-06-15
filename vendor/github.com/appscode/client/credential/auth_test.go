package credential

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var a1 = &Auth{
	Namespace: "appscode",
	Username:  "sadlil",
	Secret:    "12345",
	AuthType:  Basic,
}

var a2 = &Auth{
	Namespace: "appscode",
	Secret:    "12345",
	AuthType:  Bearer,
}

func TestBasicAuthProvider(t *testing.T) {
	bAuth := NewBasicAuth("appscode", "sadlil", "12345")
	assert.Equal(t, bAuth, a1)

	cred := bAuth.Credential()
	assert.NotNil(t, cred)
	assert.Equal(t, cred, NewBasicAuthCredential("appscode", "sadlil", "12345"))
	assert.Equal(t, reflect.TypeOf(cred), reflect.TypeOf(NewBasicAuthCredential("appscode", "sadlil", "12345")))
}

func TestBearerAuthProvider(t *testing.T) {
	bAuth := NewBearerAuth("appscode", "12345")
	assert.Equal(t, bAuth, a2)

	cred := bAuth.Credential()
	assert.NotNil(t, cred)
	assert.Equal(t, cred, NewBearerAuthCredential("appscode", "12345"))
	assert.Equal(t, reflect.TypeOf(cred), reflect.TypeOf(NewBearerAuthCredential("appscode", "12345")))

	m, err := cred.GetRequestMetadata(context.Background(), "localhost")
	assert.Nil(t, err)
	assert.Equal(t, m, map[string]string{
		"Authorization": "Bearer appscode:12345",
	})
}
