package base

const (
	ServiceName = "logging-agent"
)

const (
	OffsetFilePath = "offsets.json" // the file where offsets are persisted to restart log collection from where it left
)

// these are related to the offset persistence
const (
	DefaultOffsetUpdateIntervalMs = 2000 // default update interval for persistence of offset to the file
	DefaultOffsetRetryTimeOutMs   = 1000 // default retry timeout for persisting the offsets on failure
	DefaultOffsetRetryCount       = 1    // default retry count for persisting the offsets on failure
)

// related to the log collection from the file
const (
	DefaultLogCollectionBatchSize  = 100 // default batch size of the logs to be read from the file
	DefaultLogCollectionIntervalMs = 0   // default value of the log collection interval
)

// these are related to the log delivery endpoint
const (
	DefaultDeliveryRetryCount     = 1    // default retry timeout for calling the remote endpoint on failure
	DefaultDeliveryRetryTimeoutMs = 1000 // default retry count for calling the remote endpoint on failure
)
