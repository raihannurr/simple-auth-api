package utils_test

import (
	"errors"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestPanicIfError(t *testing.T) {
	assert.Panics(t, func() {
		utils.PanicIfError(errors.New("test error"))
	})
}
