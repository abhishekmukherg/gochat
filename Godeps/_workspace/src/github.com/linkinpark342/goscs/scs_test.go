package goscs

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type generateExp struct {
	input    []byte
	expected string
}

type fakeTimer struct {
	toReturn time.Time
}

func (f fakeTimer) Now() time.Time {
	return f.toReturn
}

func TestGenerate(t *testing.T) {
	data := []generateExp{
		{[]byte("abcdefghijklmno"), "abcdefghijklmno\x01"},
		{[]byte("abcdefghijklmn"), "abcdefghijklmn\x02\x02"},
		{[]byte("abcdefghijklmnopq"), "abcdefghijklmnopq\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f"},
		{[]byte("a"), "a\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f"},
		{[]byte(""), ""},
	}

	scs := NewMgr([]byte("deadbedwasfed123"))
	for _, exp := range data {
		cookie, err := scs.Generate(exp.input)
		if err != nil {
			t.Errorf("Got a non-nil error encoding %q", err)
		}

		result, err := scs.Parse(cookie)
		if err != nil {
			t.Errorf("Got a non-nil error decoding %q", err)
		}

		if !bytes.Equal(exp.input, result) {
			t.Errorf("Parse(Generate(%q)) = %q", exp.input, result)
		}
	}
}

func TestGenerateHasTime(t *testing.T) {
	data := []byte("abcdefjiklafehwfl")
	expected := "NTA5ZWRjNzA="
	scs := NewMgr([]byte("deadbedwasfed123"))
	scs.timer = fakeTimer{time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)}
	cookie, _ := scs.Generate(data)
	if !strings.Contains(cookie, expected) {
		t.Errorf("Generate(%q) did not contain %q", cookie, expected)
	}
}

func TestGenerateMismatch(t *testing.T) {
	data := []byte("abcdefjiklafehwfl")
	scs := NewMgr([]byte("deadbedwasfed123"))
	cookies := make([]string, 10)
	scs.timer = fakeTimer{time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)}
	for i := 0; i < 10; i++ {
		cookies[i], _ = scs.Generate(data)
	}
	for i, ci := range cookies {
		for j, cj := range cookies {
			if i == j {
				continue
			}
			if ci == cj {
				t.Errorf("%q and %q shouldn't match", cookies[i], cookies[j])
			}
		}
	}
}

func TestChangeKey(t *testing.T) {
	data := []byte("abcdefjiklafehwfl")
	scs := NewMgr([]byte("deadbedwasfed123"))
	cookie, _ := scs.Generate(data)
	scs = NewMgr([]byte("qrnqorqjnfsrq123"))
	user, err := scs.Parse(cookie)
	if user != nil {
		t.Error("Somehow managed to decrypt with a different key")
	}
	if fmt.Sprint(err) != "Bad Input" {
		t.Error("Got an unexpected error: ", err)
	}
}

func BenchmarkGenerateCookie(b *testing.B) {
	data := []byte("abcdefjiklafehwflakjhfawehfliuawefhlaiuwehlfaw")
	scs := NewMgr([]byte("deadbedwasfed123"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scs.Generate(data)
	}
}

func BenchmarkGenerateMgr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMgr([]byte("deadbedwasfed123"))
	}
}
