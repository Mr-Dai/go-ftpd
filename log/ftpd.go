package log

import "github.com/op/go-logging"

type ftpdLogger struct {
	logger *logging.Logger
}

func FtpdLogger() *ftpdLogger {
	return &ftpdLogger{logger: logger}
}

func (l *ftpdLogger) Print(sessionId string, message interface{}) {
	l.logger.Infof("%s %s", sessionId, message)
}

func (l *ftpdLogger) Printf(sessionId string, format string, v ...interface{}) {
	realV := make([]interface{}, len(v)+1)
	realV[0] = sessionId
	for i, value := range v {
		realV[i+1] = value
	}

	l.logger.Infof("%s "+format, realV...)
}

func (l *ftpdLogger) PrintCommand(sessionId string, command string, params string) {
	l.logger.Infof("%s > %s %s", sessionId, command, params)
}

func (l *ftpdLogger) PrintResponse(sessionId string, code int, message string) {
	l.logger.Infof("%s < %d %s", sessionId, code, message)
}
