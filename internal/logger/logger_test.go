package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	q2 := New("", "debug")
	q2.Trace("trace")
	q2.Warn("warn")
	q2.Error("er")
}
