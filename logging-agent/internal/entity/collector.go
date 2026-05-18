package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/krishnaZawar/distributed-logger/logger-sdk"
	"github.com/krishnaZawar/distributed-logger/logging-agent/internal/base"
	"github.com/krishnaZawar/distributed-logger/logging-agent/internal/config"
)

type LoggingAgent struct {
	mu      sync.RWMutex
	offsets map[string]int64
	readers map[string]io.ReadSeeker
	logger  *logger.Logger
}

// Creates a new logging agent to read from the readSeeker at the offset it left at
func NewLoggingAgent() *LoggingAgent {
	cfg := config.Get()
	offsets := map[string]int64{}
	file, err := os.OpenFile(base.OffsetFilePath, os.O_RDONLY, 0644)
	if err == nil {
		data, err := io.ReadAll(file)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == nil {
			err = json.Unmarshal(data, &offsets)
			if err != nil {
				panic(err)
			}
		}
	}
	file.Close()

	readers := map[string]io.ReadSeeker{}
	for _, filepath := range cfg.LogFiles {
		trackFile, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		readers[filepath] = trackFile
		_, ok := offsets[filepath]
		if !ok {
			offsets[filepath] = 0
		}
	}

	agentLogFile, err := os.OpenFile(cfg.AgentLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return &LoggingAgent{
		offsets: offsets,
		readers: readers,
		logger:  logger.New(base.ServiceName, agentLogFile),
	}
}

// Read the logs from the file from the offset it has read till
func (agent *LoggingAgent) Read() {
	for filepath, reader := range agent.readers {
		go func(filepath string, reader io.ReadSeeker) {
			for {
				agent.mu.RLock()
				offset, _ := agent.offsets[filepath]
				agent.mu.RUnlock()
				_, err := reader.Seek(offset, io.SeekStart)
				if err != nil {
					agent.logger.Error().Msgf("Error seeking file position: %s", err.Error())
					continue
				}

				scanner := json.NewDecoder(reader)

				lastOffset := offset

				logs := []log{}

				for {
					var entry log
					err := scanner.Decode(&entry)
					if err == io.EOF {
						break
					}
					if err != nil {
						agent.logger.Error().Msgf("Error Bad Json: %s", err.Error())
						break
					}

					logs = append(logs, entry)

					lastOffset = offset + scanner.InputOffset()

					if len(logs) == config.Get().LogReadBatchSize {
						break
					}
				}

				// deliver the logs
				// The offsets should not update until the logs are delivered
				// This guarantees atleast once delivery and guarantees log delivery
				err = agent.deliverLogs(logs)
				if err != nil {
					continue
				}

				agent.mu.Lock()
				agent.offsets[filepath] = lastOffset
				agent.mu.Unlock()

				time.Sleep(base.LogCollectionInterval)
			}
		}(filepath, reader)
	}
	for {
		// persist the offset periodically
		agent.persistOffset()

		time.Sleep(base.OffsetUpdateInterval)
	}
}

// persists the read offset to restart from the point the agent left on crashes.
// It also implements a retry mechanism for more durable offset persistence.
//
// If the offset does not persist even after retry, it will retry on the next read call but will not crash the agent
func (agent *LoggingAgent) persistOffset() {
	agent.mu.RLock()
	writeData, err := json.Marshal(agent.offsets)
	agent.mu.RUnlock()
	// no retry for marshalling error as it is consistent and will result in the same path
	// ideally should not occur
	if err != nil {
		return
	}
	err = agent.writeOffsetToFile(writeData)
	if err != nil {
		for i := 1; i <= base.OffsetRetryCount; i++ {
			time.Sleep(base.OffsetRetryTimeOut)
			agent.logger.Info().Msgf("Retry offset persistence count: %d", i)
			err = agent.writeOffsetToFile(writeData)
			if err == nil {
				agent.logger.Info().Msgf("Offset Persistence successful after %d retries", i)
				return
			}
		}
		agent.logger.Error().Msgf("Error: could not update offset data after %d retries: %s", base.OffsetRetryCount, err.Error())
	}
}

// writes the offset to the file
//
// returns error on failure
func (agent *LoggingAgent) writeOffsetToFile(writeData []byte) error {
	file, err := os.OpenFile(base.OffsetFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		agent.logger.Error().Msgf("Error opening offset file: %s", err.Error())
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	n, err := file.Write(writeData)
	if err == nil && n != len(writeData) {
		err = io.ErrShortWrite
		agent.logger.Error().Msgf("Error writing to offset file: %s", err.Error())
		return err
	}
	if err != nil {
		agent.logger.Error().Msgf("Error writing to offset file: %s", err.Error())
		return err
	}
	return nil
}

// sends the logs to the configured endpoint in config.yaml.
// It also implements a retry mechanism for more durable log delivery.
//
// returns an error on failure of delivery
func (agent *LoggingAgent) deliverLogs(logs []log) error {
	logData, err := json.Marshal(logs)
	if err != nil {
		agent.logger.Error().Msgf("Error: could not marshal logs: %s", err.Error())
		// no retry for marshalling error as it is consistent and will result in the same path
		// ideally should not occur
		return err
	}
	expectedStatusCode := config.Get().DeliveryDetails.ExpectedStatusCode
	resp, err := agent.callDeliveryEndpoint(logData)
	if err != nil || resp.StatusCode != expectedStatusCode {
		for i := 1; i <= base.DeliveryRetryCount; i++ {
			time.Sleep(base.DeliveryRetryTimeout)
			agent.logger.Info().Msgf("Log delivery retry count: %d", i)
			resp, err = agent.callDeliveryEndpoint(logData)
			if err == nil && resp.StatusCode == expectedStatusCode {
				agent.logger.Info().Msgf("Successfully delivered Logs after retry count: %d", i)
				return nil
			}
		}
		if err == nil {
			err = fmt.Errorf("expected status code not found, expectedStatusCode: %d, statusCode found: %d", expectedStatusCode, resp.StatusCode)
		}
		agent.logger.Error().Msgf("Error: could not deliver logs after %d retries: %s", base.DeliveryRetryCount, err.Error())
		return err
	}
	return nil
}

var client http.Client = http.Client{}

// calls the endpoint configured in config.yaml
//
// returns error on failure of delivery
func (agent *LoggingAgent) callDeliveryEndpoint(logData []byte) (*http.Response, error) {
	cfg := config.Get()
	req, err := http.NewRequest(cfg.DeliveryDetails.Method, cfg.DeliveryDetails.Endpoint, bytes.NewBuffer(logData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		agent.logger.Error().Msgf("Error: Could not form request: %s", err.Error())
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		agent.logger.Error().Msgf("%s", err.Error())
		return nil, err
	}
	return resp, nil
}
