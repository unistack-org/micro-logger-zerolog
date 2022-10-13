package zerolog // import "go.unistack.org/micro-logger-zerolog/v3"

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"go.unistack.org/micro/v3/logger"
)

type Mode uint8

const (
	Production Mode = iota
	Development
)

type zeroLogger struct {
	zLog zerolog.Logger
	opts Options
}

func (l *zeroLogger) Level(lvl logger.Level) {
}

func (l *zeroLogger) Clone(opts ...logger.Option) logger.Logger {
	nl := &zeroLogger{zLog: l.zLog, opts: l.opts}
	for _, o := range opts {
		o(&nl.opts.Options)
	}
	nl.opts.Options.Fields = append(l.opts.Options.Fields, nl.opts.Options.Fields...)
	_ = nl.Init()
	return nl
}

func (l *zeroLogger) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if zl, ok := l.opts.Context.Value(loggerKey{}).(zerolog.Logger); ok {
		l.zLog = zl
		return nil
	}

	if hs, ok := l.opts.Context.Value(hooksKey{}).([]zerolog.Hook); ok {
		l.opts.Hooks = hs
	}
	if tf, ok := l.opts.Context.Value(timeFormatKey{}).(string); ok {
		l.opts.TimeFormat = tf
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if useDefault, ok := l.opts.Context.Value(useAsDefaultKey{}).(bool); ok && useDefault {
		l.opts.UseAsDefault = useDefault
	}
	if devMode, ok := l.opts.Context.Value(developmentModeKey{}).(bool); ok && devMode {
		l.opts.Mode = Development
	}
	if prodMode, ok := l.opts.Context.Value(productionModeKey{}).(bool); ok && prodMode {
		l.opts.Mode = Production
	}

	// RESET
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = nil
	zerolog.CallerSkipFrameCount = 4

	switch l.opts.Mode {
	case Development:
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		consOut := zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				if len(l.opts.TimeFormat) > 0 {
					w.TimeFormat = l.opts.TimeFormat
				}
				w.Out = l.opts.Out
				w.NoColor = false
			},
		)
		// level = logger.DebugLevel
		l.zLog = zerolog.New(consOut).
			Level(zerolog.DebugLevel).
			With().Timestamp().Stack().Logger()
	default: // Production
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		l.zLog = zerolog.New(l.opts.Out).
			Level(zerolog.InfoLevel).
			With().Timestamp().Stack().Logger()
	}

	// Set log Level if not default
	if l.opts.Level != 100 {
		zerolog.SetGlobalLevel(loggerToZerologLevel(l.opts.Level))
		l.zLog = l.zLog.Level(loggerToZerologLevel(l.opts.Level))
	}

	// Reporting caller
	if l.opts.ReportCaller {
		l.zLog = l.zLog.With().Caller().Logger()
	}

	// Adding hooks if exist
	for _, hook := range l.opts.Hooks {
		l.zLog = l.zLog.Hook(hook)
	}

	// Setting timeFormat
	if len(l.opts.TimeFormat) > 0 {
		zerolog.TimeFieldFormat = l.opts.TimeFormat
	}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := make(map[string]interface{}, len(l.opts.Fields)/2)
		for i := 0; i < len(l.opts.Fields); i += 2 {
			fld, ok := l.opts.Fields[i].(string)
			if !ok {
				return fmt.Errorf("field type invalid %T not string", l.opts.Fields[i])
			}
			data[fld] = l.opts.Fields[i+1]
		}
		l.zLog = l.zLog.With().Fields(data).Logger()
	}

	// Also set it as zerolog's Default logger
	if l.opts.UseAsDefault {
		zlog.Logger = l.zLog
	}

	return nil
}

func (l *zeroLogger) Fields(fields ...interface{}) logger.Logger {
	return l.Clone(logger.WithFields(fields...))
}

func (l *zeroLogger) V(level logger.Level) bool {
	return l.zLog.GetLevel() <= loggerToZerologLevel(level)
}

func (l *zeroLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

func (l *zeroLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

func (l *zeroLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

func (l *zeroLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

func (l *zeroLogger) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.TraceLevel, args...)
}

func (l *zeroLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	l.opts.ExitFunc(1)
}

func (l *zeroLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.InfoLevel, msg, args...)
}

func (l *zeroLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.ErrorLevel, msg, args...)
}

func (l *zeroLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.WarnLevel, msg, args...)
}

func (l *zeroLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.DebugLevel, msg, args...)
}

func (l *zeroLogger) Tracef(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.TraceLevel, msg, args...)
}

func (l *zeroLogger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.FatalLevel, msg, args...)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	l.opts.ExitFunc(1)
}

func (l *zeroLogger) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	if !l.V(level) {
		return
	}

	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(loggerToZerologLevel(level)).Msg(msg)
}

func (l *zeroLogger) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.zLog.WithLevel(loggerToZerologLevel(level)).Msgf(format, args...)
}

func (l *zeroLogger) String() string {
	return "zerolog"
}

func (l *zeroLogger) Options() logger.Options {
	return l.opts.Options
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...logger.Option) logger.Logger {
	// Default options
	options := Options{
		Options:      logger.NewOptions(opts...),
		ReportCaller: false,
		UseAsDefault: false,
		Mode:         Production,
		ExitFunc:     os.Exit,
	}

	l := &zeroLogger{opts: options}
	return l
}

func loggerToZerologLevel(level logger.Level) zerolog.Level {
	switch level {
	case logger.TraceLevel:
		return zerolog.TraceLevel
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	case logger.FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func ZerologToLoggerLevel(level zerolog.Level) logger.Level {
	switch level {
	case zerolog.TraceLevel:
		return logger.TraceLevel
	case zerolog.DebugLevel:
		return logger.DebugLevel
	case zerolog.InfoLevel:
		return logger.InfoLevel
	case zerolog.WarnLevel:
		return logger.WarnLevel
	case zerolog.ErrorLevel:
		return logger.ErrorLevel
	case zerolog.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
