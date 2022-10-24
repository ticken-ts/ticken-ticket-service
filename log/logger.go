package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"ticken-ticket-service/env"
)

const LogFilename = "logs"

var TickenLogger zerolog.Logger

func InitGlobalLogger() zerolog.Logger {
	TickenLogger = zerolog.New(getOutput()).Level(getLogLevel()).With().Timestamp().Logger()
	return TickenLogger
}

func getLogLevel() zerolog.Level {
	switch env.TickenEnv.Env {
	case env.ProdEnv:
		return zerolog.InfoLevel
	case env.StageEnv:
		return zerolog.DebugLevel
	default: // test and dev
		return zerolog.TraceLevel
	}
}

func getOutput() io.Writer {
	switch env.TickenEnv.Env {
	case env.ProdEnv:
		return buildFileWriter()
	case env.StageEnv:
		return buildConsoleWriter()
	default: // test and dev
		return buildConsoleWriter()
	}
}

func buildConsoleWriter() io.Writer {
	consoleWriter := zerolog.NewConsoleWriter()

	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	return consoleWriter
}

func buildFileWriter() io.Writer {
	// if not exists, it is created
	logs, err := os.Create(LogFilename)
	if err != nil {
		panic(err)
	}

	return logs
}
