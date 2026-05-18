package base

import "time"

const (
	ServiceName = "logging-agent"
)

const (
	OffsetUpdateInterval = 5 * time.Second
	OffsetFilePath       = "offsets.json"
	OffsetRetryTimeOut   = time.Second
	OffsetRetryCount     = 1
)

const (
	LogCollectionInterval = 2 * time.Second
)

const (
	DeliveryRetryCount   = 1
	DeliveryRetryTimeout = time.Second
)
