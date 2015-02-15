package goscs

import (
	"bytes"
	crand "crypto/rand"
	"math/rand"
	"testing"
)

type expectations struct {
	input, expected []byte
}

func TestAddPadding(t *testing.T) {
	data := []expectations{
		{[]byte("abcdefghijklmno"), []byte("abcdefghijklmno\x01")},
		{[]byte("abcdefghijklmn"), []byte("abcdefghijklmn\x02\x02")},
		{[]byte("abcdefghijklmnopq"),
			[]byte("abcdefghijklmnopq\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f")},
		{[]byte("a"), []byte("a\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f")},
		{[]byte(""), bytes.Repeat([]byte{byte(0x10)}, 0x10)},
	}

	for _, exp := range data {
		padded := addPadding(exp.input)
		if !bytes.Equal(padded, exp.expected) {
			t.Errorf("addPadding(%q) = %q, expected %q", exp.input, padded, exp.expected)
		}
	}
}

func BenchmarkAddPaddingAlphabet(b *testing.B) {
	src := []byte("abcdefghijklmnopqrstuvwxyz")
	length := len(src)

	for n := 0; n < b.N; n++ {
		for i := 0; i < length; i++ {
			addPadding(src[:i])
		}
	}
}

func BenchmarkAddPaddingFuzz(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fuzzSize := rand.Uint32() % (1 << 15)
		fuzz := make([]byte, fuzzSize)
		crand.Read(fuzz)
		addPadding(fuzz)
	}
}

func BenchmarkRemovePaddingAlphabet(b *testing.B) {
	rawString := []byte("abcdefghijklmnopqrstuvwxyz")
	length := len(rawString)
	src := make([][]byte, len(rawString))
	for i := 0; i < length; i++ {
		src[i] = addPadding(rawString[:i])
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < length; i++ {
			removePadding(src[i])
		}
	}
}

func BenchmarkAddRemovePaddingFuzz(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fuzzSize := rand.Uint32() % (1 << 15)
		fuzz := make([]byte, fuzzSize)
		crand.Read(fuzz)
		removePadding(addPadding(fuzz))
	}
}
