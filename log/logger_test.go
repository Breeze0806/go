package log

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
)

type mockWriter struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

func newMockWriter(buf *bytes.Buffer) *mockWriter {
	return &mockWriter{
		buf: buf,
	}
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buf.Write(p)
}

func TestNewDefaultLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	SetLogger(NewDefaultLogger(newMockWriter(buf), DebugLevel, "datax"))
	testCases := []struct {
		printf func(string, ...interface{})
		format string
		args   []interface{}
	}{
		{
			printf: GetLogger().Debugf,
			format: "debug %d",
			args:   []interface{}{DebugLevel},
		},
		{
			printf: GetLogger().Infof,
			format: "info %d",
			args:   []interface{}{InfoLevel},
		},
		{
			printf: GetLogger().Errorf,
			format: "error %d",
			args:   []interface{}{ErrorLevel},
		},
	}

	for _, v := range testCases {
		buf.Reset()
		v.printf(v.format, v.args...)
		a := strings.Split(buf.String(), ": ")
		out := a[len(a)-1]
		out = out[:len(out)-1]
		want := fmt.Sprintf(v.format, v.args...)

		if want != out {
			t.Fatalf("want != out want: %v[%v] out: %v[%v] log: %v.", want, len(want), out, len(out), buf.String())
		}
	}
}

func TestDefaultLogger_Print(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	SetLogger(NewDefaultLogger(newMockWriter(buf), DebugLevel, "datax"))
	testCases := []struct {
		print func(...interface{})
		args  []interface{}
	}{
		{
			print: GetLogger().Print,
			args:  []interface{}{DebugLevel},
		},
	}

	for _, v := range testCases {
		buf.Reset()
		v.print(v.args...)
		a := strings.Split(buf.String(), ": ")
		out := a[len(a)-1]
		out = out[:len(out)-1]
		want := fmt.Sprint(v.args...)

		if want != out {
			t.Fatalf("want != out want: %v[%v] out: %v[%v] log: %v.", want, len(want), out, len(out), buf.String())
		}
	}
}

func Test_defaultLogger_Printf(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	SetLogger(NewDefaultLogger(newMockWriter(buf), DebugLevel, "datax"))
	testCases := []struct {
		printf func(string, ...interface{})
		format string
		args   []interface{}
	}{
		{
			printf: GetLogger().Printf,
			format: "debug %d",
			args:   []interface{}{DebugLevel},
		},
		{
			printf: GetLogger().Printf,
			format: "info %d",
			args:   []interface{}{InfoLevel},
		},
		{
			printf: GetLogger().Printf,
			format: "error %d",
			args:   []interface{}{ErrorLevel},
		},
	}

	for _, v := range testCases {
		buf.Reset()
		v.printf(v.format, v.args...)
		a := strings.Split(buf.String(), ": ")
		out := a[len(a)-1]
		out = out[:len(out)-1]
		want := fmt.Sprintf(v.format, v.args...)

		if want != out {
			t.Fatalf("want != out want: %v[%v] out: %v[%v] log: %v.", want, len(want), out, len(out), buf.String())
		}
	}
}

func TestRegisterInitFuncs(t *testing.T) {
	var log1 Logger
	var log2 Logger
	f1 := func() {
		log1 = GetLogger()
	}
	f2 := func() {
		log2 = GetLogger()
	}
	RegisterInitFuncs(f1)
	RegisterInitFuncs(f2)
	buf := bytes.NewBuffer(nil)
	SetLogger(NewDefaultLogger(newMockWriter(buf), DebugLevel, "datax"))
	if log1 != GetLogger() {
		t.Fatalf("want != out want: %p out: %p", GetLogger(), log1)
	}
	if log2 != GetLogger() {
		t.Fatalf("want != out want: %p out: %p", GetLogger(), log2)
	}
}
