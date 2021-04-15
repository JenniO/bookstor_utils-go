package rest_errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError("this is the message", errors.New("database error"))
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.ErrStatus)
	assert.EqualValues(t, "this is the message", err.ErrMessage)
	assert.EqualValues(t, "internal_server_error", err.Error)

	assert.NotNil(t, err.ErrCauses)
	assert.EqualValues(t, 1, len(err.ErrCauses))
	assert.EqualValues(t, "database error", err.ErrCauses[0])
}

func TestNewBadRequestError(t *testing.T) {

}

func TestNewNotFoundError(t *testing.T) {

}

func TestNewError(t *testing.T) {

}
