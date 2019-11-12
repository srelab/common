package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"

	"github.com/srelab/common/random"

	"github.com/sirupsen/logrus"
	"github.com/srelab/common/file"
)

// Level describes the log severity level.
type Level uint8

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type Config struct {
	File  string `yaml:"File"`
	Level string `yaml:"Level"`
}

// Logger is an interface that describes logging.
type Logger interface {
	With(key string, value interface{}) Logger
	WithError(err error) Logger

	SetLevel(level Level)
	SetOut(out io.Writer)

	Trace(...interface{})
	Debug(...interface{})
	Print(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})

	Tracef(string, ...interface{})
	Debugf(string, ...interface{})
	Printf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})

	Traceln(...interface{})
	Debugln(...interface{})
	Println(...interface{})
	Infoln(...interface{})
	Warnln(...interface{})
	Errorln(...interface{})
	Fatalln(...interface{})
	Panicln(...interface{})
}

type logger struct {
	entry *logrus.Entry
}

// With attaches a key-value pair to a logger.
func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// WithError attaches an error to a logger.
func (l logger) WithError(err error) Logger {
	return logger{l.entry.WithError(err)}
}

// SetLevel sets the level of a logger.
func (l logger) SetLevel(level Level) {
	l.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination for a logger.
func (l logger) SetOut(out io.Writer) {
	l.entry.Logger.Out = out
}

// Trace logs a message at level Trace on the standard logger.
func (l logger) Trace(args ...interface{}) {
	l.sourced().Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Print logs a message at level Print on the standard logger.
func (l logger) Print(args ...interface{}) {
	l.sourced().Print(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Error logs a message at Error Info on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Panic logs a message at level Panic on the standard logger.
func (l logger) Panic(args ...interface{}) {
	l.sourced().Panic(args...)
}

func (l logger) Tracef(format string, args ...interface{}) {
	l.sourced().Tracef(format, args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

func (l logger) Printf(format string, args ...interface{}) {
	l.sourced().Printf(format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

func (l logger) Panicf(format string, args ...interface{}) {
	l.sourced().Panicf(format, args...)
}

func (l logger) Traceln(args ...interface{}) {
	l.sourced().Traceln(args...)
}

func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

func (l logger) Println(args ...interface{}) {
	l.sourced().Println(args...)
}

func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

func (l logger) Panicln(args ...interface{}) {
	l.sourced().Panicln(args...)
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	_, _file, line, ok := runtime.Caller(2)

	if !ok {
		_file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(_file, "/")
		_file = _file[slash+1:]
	}

	return l.entry.WithField("src", fmt.Sprintf("%s:%d", _file, line))
}

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// New returns a new logger.
func New() Logger {
	return logger{entry: logrus.NewEntry(origLogger)}
}

// Base returns the base logger.
func Base() Logger {
	return baseLogger
}

// Initialize the logger with config
// When path is not legal, the current path will be used.
// Multiwriter by default
func Init(config Config) {
	fp, fn := file.Dir(config.File), file.Basename(config.File)
	if err := file.EnsureDirRW(fp); err != nil || (file.IsExist(config.File) && !file.IsFile(config.File))  {
		fp, fn = "./", random.New().String(8, random.Lowercase)+".log"
	}

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	SetLevel(Level(level))
	SetOut(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   path.Join(fp, fn),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		LocalTime:  true,
	}))
}

// SetLevel sets the Level of the base logger
func SetLevel(level Level) {
	baseLogger.entry.Logger.Level = logrus.Level(level)
}

// GetLevel gets the level of a logger.
func GetLevel() Level {
	return Level(baseLogger.entry.Logger.Level)
}

// SetOut sets the output destination base logger
func SetOut(out io.Writer) {
	baseLogger.entry.Logger.Out = out
}

func With(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

func WithError(err error) Logger {
	return logger{entry: baseLogger.sourced().WithError(err)}
}

func Trace(args ...interface{}) {
	baseLogger.sourced().Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	baseLogger.sourced().Tracef(format, args...)
}

func Traceln(args ...interface{}) {
	baseLogger.sourced().Traceln(args...)
}

func Debug(args ...interface{}) {
	baseLogger.sourced().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	baseLogger.sourced().Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

func Print(args ...interface{}) {
	baseLogger.sourced().Print(args...)
}

func Printf(format string, args ...interface{}) {
	baseLogger.sourced().Printf(format, args...)
}

func Println(args ...interface{}) {
	baseLogger.sourced().Println(args...)
}

func Info(args ...interface{}) {
	baseLogger.sourced().Info(args...)
}

func Infof(format string, args ...interface{}) {
	baseLogger.sourced().Infof(format, args...)
}

func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

func Warn(args ...interface{}) {
	baseLogger.sourced().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	baseLogger.sourced().Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

func Error(args ...interface{}) {
	baseLogger.sourced().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	baseLogger.sourced().Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

func Fatal(args ...interface{}) {
	baseLogger.sourced().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced().Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

func Panic(args ...interface{}) {
	baseLogger.sourced().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	baseLogger.sourced().Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	baseLogger.sourced().Panicln(args...)
}
