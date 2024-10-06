package helper

import "testing"

func TestHelloWorld(t *testing.T) {
	result := HelloWorld("eka")
	if result != "hello eka" {
		t.FailNow()
	}
}
