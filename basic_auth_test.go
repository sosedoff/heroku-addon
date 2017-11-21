package addon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseBasicAuth(t *testing.T) {
	_, err := parseBasicAuth("foobar")
	assert.NotNil(t, err)
	assert.Error(t, ErrNoBasicAuth, err.Error())

	_, err = parseBasicAuth("Basic qwe123")
	assert.NotNil(t, err)
	assert.Error(t, ErrNoBasicAuth, err.Error())

	auth, err := parseBasicAuth("Basic QWxhZGRpbjpPcGVuU2VzYW1l")
	assert.NoError(t, err)
	assert.Equal(t, "Aladdin", auth.Username)
	assert.Equal(t, "OpenSesame", auth.Password)
}
