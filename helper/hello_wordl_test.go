package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	t.Skip("skip unit test")
	result := HelloWorld("eka")
	if result != "hello eka" {
		t.Fatal("tidak sesuai")
	}
}

func TestHelloWorldAssert(t *testing.T) {
	result := HelloWorld("eka")
	assert.Equal(t, "hello eka", result)
}
