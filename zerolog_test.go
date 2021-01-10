package zero

import (
	"os"
	"testing"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "zerolog" {
		t.Errorf("error: name expected 'zerolog' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

// func ExampleWithOut() {
// 	l := NewLogger(WithOut(os.Stdout), WithProductionMode())

// 	l.Logf(logger.InfoLevel, "testing: %s", "logf")

// 	// Output:
// 	// {"level":"info","time":"2020-02-13T20:55:24-08:00","message":"testing: logf"}
// }

func TestSetLevel(t *testing.T) {
	l := NewLogger()

	l.SetLevel(logger.DebugLevel)
	l.Logf(logger.DebugLevel, "test show debug: %s", "debug msg")

	l.SetLevel(logger.InfoLevel)
	l.Logf(logger.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	l := NewLogger(ReportCaller())

	l.Logf(logger.InfoLevel, "testing: %s", "WithReportCaller")
}
func TestWithOut(t *testing.T) {
	l := NewLogger(WithOut(os.Stdout))

	l.Logf(logger.InfoLevel, "testing: %s", "WithOut")
}

func TestWithPretty(t *testing.T) {
	l := NewLogger(WithDevelopmentMode())

	l.Logf(logger.InfoLevel, "testing: %s", "WithPretty")
}
func TestWithLevelFieldName(t *testing.T) {
	l := NewLogger(WithGCPMode())

	l.Logf(logger.InfoLevel, "testing: %s", "WithLevelFieldName")
	// reset `LevelFieldName` to make other tests pass.
	NewLogger(WithProductionMode())
}

func TestWithFields(t *testing.T) {
	l := NewLogger()

	l.Fields([]logger.Field{
		{
			Key:   "sumo",
			Type:  logger.StringType,
			Value: "demo",
		},
		{
			Key:   "human",
			Type:  logger.BoolType,
			Value: true,
		},
		{
			Key:   "age",
			Type:  logger.Int32Type,
			Value: 99,
		},
	}...).Logf(logger.InfoLevel, "testing: %s", "WithFields")
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	l := NewLogger(WithHooks([]zerolog.Hook{simpleHook}))

	l.Logf(logger.InfoLevel, "testing: %s", "WithHooks")
}
