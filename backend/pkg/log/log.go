package log

import (
	"strings"

	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var Module = fx.Options(
	fx.Provide(ProvideLogger),
	// inject custom logger for fx
	fx.WithLogger(func(l *Logger) fxevent.Logger {
		return l
	}),
)

type Logger struct {
	*logrus.Logger
}

func ProvideLogger(config config.Config) *Logger {
	logger := logrus.New()

	if config.Environment.IsLocal() {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return &Logger{
		Logger: logger,
	}
}

// LogEvent handles a log event for fx application container
// taken from https://github.com/takt-corp/fx-logrus/blob/main/logrus.go
func (l *Logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.WithFields(logrus.Fields{
			"callee": e.FunctionName,
			"caller": e.CallerName,
		}).Debug("on start hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"callee": e.FunctionName,
				"caller": e.CallerName,
			}).Errorf("on start hook failed: %v", e.Err)
		} else {
			l.Logger.WithFields(logrus.Fields{
				"callee":  e.FunctionName,
				"caller":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Debug("on start hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.WithFields(logrus.Fields{
			"callee": e.FunctionName,
			"caller": e.CallerName,
		}).Debug("on stop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"callee": e.FunctionName,
				"caller": e.CallerName,
			}).Errorf("on stop hook failed: %v", e.Err)
		} else {
			l.Logger.WithFields(logrus.Fields{
				"callee":  e.FunctionName,
				"caller":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Debug("on stop hook executed")
		}
	case *fxevent.Supplied:
		l.Logger.WithFields(logrus.Fields{
			"type":   e.TypeName,
			"module": e.ModuleName,
		}).Debugf("supplied: %v", e.Err)
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.WithFields(logrus.Fields{
				"constructor": e.ConstructorName,
				"module":      e.ModuleName,
				"type":        rtype,
			}).Debug("provided")
		}
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"module": e.ModuleName,
			}).Errorf("error encountered while applying options: %v", e.Err)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.WithFields(logrus.Fields{
				"module": e.ModuleName,
				"type":   rtype,
			}).Debug("replaced")
		}
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"module": e.ModuleName,
			}).Errorf("error encountered while replacing: %v", e.Err)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.WithFields(logrus.Fields{
				"module": e.ModuleName,
				"type":   rtype,
			}).Debug("decorated")
		}
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"module": e.ModuleName,
			}).Errorf("error encountered while applying options: %v", e.Err)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Logger.WithFields(logrus.Fields{
			"function": e.FunctionName,
			"module":   e.ModuleName,
		}).Debug("invoking")
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.WithFields(logrus.Fields{
				"stack":    e.Trace,
				"function": e.FunctionName,
				"module":   e.ModuleName,
			}).Errorf("invoke failed: %v", e.Err)
		}
	case *fxevent.Stopping:
		l.Logger.Debugf("received signal: %s", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Errorf("received signal: %v", e.Err)
		}
	case *fxevent.RollingBack:
		l.Logger.Errorf("start failed, rolling back: %v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Errorf("rollback failed: %v", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Errorf("start failed: %v", e.Err)
		} else {
			l.Logger.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Errorf("custom logger initialization failed: %v", e.Err)
		} else {
			l.Logger.WithFields(logrus.Fields{
				"function": e.ConstructorName,
			}).Debug("initialized custom fxevent.Logger")
		}
	}
}
