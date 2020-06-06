package logging

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestLogInfo(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	tm := "Hello World!"
	Information(tm)

	exp := formater(LogLevelInformation, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogInfof(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	tm := "Hello World!"
	Informationf(tm)

	exp := formater(LogLevelInformation, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogWarning(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	tm := "Hello World!"
	Warning(tm)

	exp := formater(LogLevelWarning, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogWarningf(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	tm := "Hello World!"
	Warningf(tm)

	exp := formater(LogLevelWarning, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogWarningDisabled(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o).SetLogLevel(LogLevelInformation)

	tm := "Hello World!"
	Warning(tm)

	if sb.String() != "" {
		t.Errorf("expected '' but got: '%s'", sb.String())
	}
}

func TestLogDebug(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o).SetLogLevel(LogLevelDebug)

	tm := "Hello World!"
	Debug(tm)

	exp := formater(LogLevelDebug, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogDebugf(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o).SetLogLevel(LogLevelDebug)

	tm := "Hello World!"
	Debugf(tm)

	exp := formater(LogLevelDebug, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogDebugDisabled(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o).SetLogLevel(LogLevelInformation)

	tm := "Hello World!"
	Debug(tm)

	if sb.String() != "" {
		t.Errorf("expected '' but got: '%s'", sb.String())
	}
}

func TestLogError(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	te := errors.New("Hello World")
	Error(te)

	exp := formater(LogLevelError, fmt.Sprintf("%v", te))
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestLogErrorf(t *testing.T) {
	o, sb := buildTestOutput()
	ResetOutputs().WithOutput(o)

	tm := "Hello World!"
	Errorf(tm)

	exp := formater(LogLevelError, tm)
	act := sb.String()

	if exp != act {
		t.Errorf("expected '%s' but got: '%s'", exp, act)
	}
}

func TestWithOutput(t *testing.T) {
	so := StdOut()
	ResetOutputs()
	WithOutput(so)

	if len(defaultLogger.outputs) != 1 {
		t.Errorf("expected 1 output but got: %d", len(defaultLogger.outputs))
	} else {
		if defaultLogger.outputs[0] != so {
			t.Errorf("unexpected output option")
		}
	}
}

func TestSetLogLevel(t *testing.T) {
	l := LogLevelError
	SetLogLevel(l)

	if l != defaultLogger.lvl {
		t.Errorf("log level expected to be %d but got: %d", l, defaultLogger.lvl)
	}
}

func buildTestOutput() (Output, *strings.Builder) {
	sb := new(strings.Builder)
	o := &testOutput{sb}
	return o, sb
}

type testOutput struct {
	sb *strings.Builder
}

func (o *testOutput) Write(message string) {
	o.sb.WriteString(message)
}
