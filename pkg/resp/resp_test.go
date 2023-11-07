package resp

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRespSuccess(t *testing.T) {
	r := Success(nil)
	assert.Equal(t, 200, r.Code)
}

func TestClientFail(t *testing.T) {
	r := ClientFail("")
	assert.Equal(t, -1, r.Code)
}

func TestClientFailWithCode(t *testing.T) {
	r := ClientFailWithCode(1, "123")
	assert.Equal(t, 1, r.Code)
	assert.Equal(t, "123", r.Msg)
}

func TestServerFail(t *testing.T) {
	r := ServerFail("")
	assert.Equal(t, -1, r.Code)
}

func TestServerFailWithCode(t *testing.T) {
	r := ServerFailWithCode(1, "123")
	assert.Equal(t, 1, r.Code)
	assert.Equal(t, "123", r.Msg)
}
