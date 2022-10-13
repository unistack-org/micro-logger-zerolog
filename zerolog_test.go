package zerolog

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.unistack.org/micro/v3/logger"
)

func TestFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Fields("key", "val")
	nl.Infof(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
	buf.Reset()
	mnl := nl.Fields("key1", "val1")
	mnl.Infof(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key1":"val1"`)) || !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestOutput(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Infof(context.TODO(), "test logger name: %s", "name")
	if !bytes.Contains(buf.Bytes(), []byte(`test logger name`)) {
		t.Fatalf("log not redirected: %s", buf.Bytes())
	}
}

func TestName(t *testing.T) {
	l := NewLogger()
	l.Init()
	if l.String() != "zerolog" {
		t.Errorf("error: name expected 'zerolog' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithOutput(t *testing.T) {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))
	logger.DefaultLogger.Init()
	logger.Infof(context.TODO(), "testing: %s", "WithOutput")
}

func TestSetLevel(t *testing.T) {
	logger.DefaultLogger = NewLogger()
	logger.Init(logger.WithLevel(logger.DebugLevel))
	logger.Debugf(context.TODO(), "test show debug: %s", "debug msg")

	logger.Init(logger.WithLevel(logger.InfoLevel))
	logger.Debugf(context.TODO(), "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logger.DefaultLogger = NewLogger(ReportCaller())
	logger.DefaultLogger.Init()
	logger.Infof(context.TODO(), "testing: %s", "WithReportCaller")
}

func TestWithOut(t *testing.T) {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))
	logger.DefaultLogger.Init()
	logger.Infof(context.TODO(), "testing: %s", "WithOut")
}

func TestWithDevelopmentMode(t *testing.T) {
	logger.DefaultLogger = NewLogger(WithDevelopmentMode(), WithTimeFormat(time.Kitchen))
	logger.DefaultLogger.Init()
	logger.Infof(context.TODO(), "testing: %s", "DevelopmentMode")
}

func TestWithFields(t *testing.T) {
	logger.DefaultLogger = NewLogger()
	logger.DefaultLogger.Init()
	logger.Fields("sumo", "demo", "human", true, "age", 99).Infof(context.TODO(), "testing: %s", "WithFields")
}

func TestWithError(t *testing.T) {
	logger.DefaultLogger = NewLogger()
	logger.DefaultLogger.Init()
	logger.Fields("error", errors.New("I am Error")).Errorf(context.TODO(), "testing: %s", "WithError")
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	logger.DefaultLogger = NewLogger(WithHooks([]zerolog.Hook{simpleHook}))
	logger.DefaultLogger.Init()
	logger.Infof(context.TODO(), "testing: %s", "WithHooks")
}
