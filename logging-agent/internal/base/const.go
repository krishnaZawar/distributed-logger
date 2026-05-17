package base

import "time"

const (
	ServiceName = "logging-agent"
)

const (
	OffsetFilePath     = "offsets.json"
	OffsetRetryTimeOut = time.Second
)

const (
	LogReadBatchSize = 50
)
