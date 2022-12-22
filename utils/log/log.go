package log

import (
	"context"

	"github.com/ahmadmirdas/julo-test/utils/activity"

	"github.com/sirupsen/logrus"
)

const (
	MAX_LOG_ENTRY_SIZE = 8 * 1024
)

func WithContext(ctx context.Context) *logrus.Entry {
	fields := activity.GetFields(ctx)
	return logrus.WithFields(fields)
}
