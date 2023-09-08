package log

import (
	"context"

	"github.com/sirupsen/logrus"

	"go-poc/utils/activity"
)

const (
	MAX_LOG_ENTRY_SIZE = 8 * 1024
)

func WithContext(ctx context.Context) *logrus.Entry {
	fields := activity.GetFields(ctx)

	return logrus.WithFields(fields)
}

func LogOrmer(obj interface{}, prefix string) {
	logrus.Debug(prefix, obj)
}
