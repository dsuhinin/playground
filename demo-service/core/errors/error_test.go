package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Common testing constants.
//
const (
	emptyMsg = ""
	errMsg   = "error message"
	infoMsg  = "info msg"
)

//
// New :: with an empty message :: returns an empty message
//
func TestNewErrorWithoutAMessage(t *testing.T) {

	assert.Equal(t, emptyMsg, New(emptyMsg).Error())
}

//
// New :: with a message :: returns a message
//
func TestNewError(t *testing.T) {

	assert.Equal(t, errMsg, New(errMsg).Error())
}

//
// New :: with formatted error message :: returns a message
//
func TestNewWithFormattedErrorMessage(t *testing.T) {

	assert.Equal(t, errMsg, New("error %s", "message").Error())
}

//
// WithInfo :: for an empty info :: returns original error message.
//
func TestWithContextForAnEmptyInfo(t *testing.T) {
	err := New(errMsg)
	chain := err.WithMessage("")

	assert.Equal(t, errMsg, chain.Error())
}

//
// WithInfo :: for a non-empty info :: returns original error message and info.
//
func TestWithContextForANotEmptyInfo(t *testing.T) {
	err := New(errMsg)
	chain := err.WithMessage(infoMsg)
	expectedMsg := fmt.Sprintf("%s : %s", infoMsg, errMsg)

	assert.Equal(t, expectedMsg, chain.Error())
}
