package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        time.StampMilli, // the "time" field
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(frame.File), frame.Line)
		},
	}
	logrus.SetFormatter(formatter)
}
func formatFilePath(path string) string {
	array := strings.Split(path, "/")
	return array[len(array)-1]
}
