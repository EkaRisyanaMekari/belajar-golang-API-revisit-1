package helper

import (
	"fmt"
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

func TestHelloWorldMulti(t *testing.T) {
	t.Run("Case1", func(t *testing.T) {
		result := HelloWorld("eka")
		assert.Equal(t, "hello eka", result)
	})
	t.Run("Case2", func(t *testing.T) {
		result := HelloWorld("joko")
		assert.Equal(t, "hello joko", result)
	})
}

func TestTableHelloWorld(t *testing.T) {
	lists := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Eka",
			param:    "Eka",
			expected: "hello Eka",
		},
		{
			name:     "Risyana",
			param:    "Risyana",
			expected: "hello Risyana",
		},
		{
			name:     "Pribadi",
			param:    "Pribadi",
			expected: "hello Pribadi",
		},
	}

	for _, list := range lists {
		t.Run(list.name, func(t *testing.T) {
			result := HelloWorld(list.param)
			assert.Equal(t, list.expected, result)
		})
	}
}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HelloWorld("Eka")
	}
}
func BenchmarkHelloWorldMulti(b *testing.B) {
	b.Run("Benchmark1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HelloWorld("Eka")
		}
	})
	b.Run("Benchmark2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HelloWorld("Risyana")
		}
	})
}

func TestMain(m *testing.M) {
	fmt.Println("before test")
	m.Run()
	fmt.Println("after test")
}
