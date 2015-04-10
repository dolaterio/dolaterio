package dolaterio

import (
	"strings"
	"testing"
	"time"
)

var (
	testContainerEngine = &ContainerEngine{}
)

func init() {
	err := testContainerEngine.Connect()
	if err != nil {
		panic(err)
	}
}

func assertString(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("Expected \"%s\", got \"%v\"", s1, s2)
	}
}
func assertStringContains(t *testing.T, s1, s2 string) {
	if strings.Contains(s1, s2) {
		t.Errorf("Expected \"%s\" to contain \"%v\"", s2, s1)
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("Expected \"%v\" to be nil", v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Errorf("Expected not to be nil")
	}
}

func assertMaxDuration(t *testing.T, d1, d2 time.Duration) {
	if d2 > d1 {
		t.Errorf("Expected %0.4fs to be shorter than %0.4fs", d2.Seconds(), d1.Seconds())
	}
}
