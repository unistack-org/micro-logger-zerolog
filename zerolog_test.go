package zerolog

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unistack-org/micro/v3/logger"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "zerolog" {
		t.Errorf("error: name expected 'zerolog' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithOutput(t *testing.T) {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))

	logger.Infof("testing: %s", "WithOutput")
}

func TestSetLevel(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	logger.Init(logger.WithLevel(logger.DebugLevel))
	logger.Debugf("test show debug: %s", "debug msg")

	logger.Init(logger.WithLevel(logger.InfoLevel))
	logger.Debugf("test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logger.DefaultLogger = NewLogger(ReportCaller())

	logger.Infof("testing: %s", "WithReportCaller")
}

func TestWithOut(t *testing.T) {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))

	logger.Infof("testing: %s", "WithOut")
}

func TestWithDevelopmentMode(t *testing.T) {
	logger.DefaultLogger = NewLogger(WithDevelopmentMode(), WithTimeFormat(time.Kitchen))

	logger.Infof("testing: %s", "DevelopmentMode")
}

func TestWithFields(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	logger.Fields(map[string]interface{}{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	}).Infof("testing: %s", "WithFields")
}

func TestWithError(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	logger.Fields(map[string]interface{}{"error": errors.New("I am Error")}).Errorf("testing: %s", "WithError")
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	logger.DefaultLogger = NewLogger(WithHooks([]zerolog.Hook{simpleHook}))

	logger.Infof("testing: %s", "WithHooks")
}
