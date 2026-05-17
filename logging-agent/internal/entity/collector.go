package entity

import (
	"encoding/json"
	"fmt"
	"io"
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
			fmt.Println(string(data))
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

				logs := []Log{}

				for {
					var entry Log
					err := scanner.Decode(&entry)
					if err == io.EOF {
						break
					}
					if err != nil {
						agent.logger.Error().Msgf("Error Bad Json: %s", err.Error())
						break
					}

					logs = append(logs, entry)

					lastOffset = scanner.InputOffset()

					if len(logs) == base.LogReadBatchSize {
						break
					}
				}

				// call delivery endpoint
				fmt.Println(logs)

				agent.mu.Lock()
				agent.offsets[filepath] = lastOffset
				agent.mu.Unlock()

				time.Sleep(2 * time.Second)
			}
		}(filepath, reader)
	}
	for {
		// persist the offset
		agent.persistOffset()

		time.Sleep(2 * time.Second)
	}
}

// persists the read offset to restart from the point the agent left on crashes.
// It also implements a single retry mechanism for more durable offset persistence.
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
		time.Sleep(base.OffsetRetryTimeOut)
		agent.logger.Info().Msg("Retry offset persistence for durability")
		err = agent.writeOffsetToFile(writeData)
		if err != nil {
			agent.logger.Error().Msgf("Error: could not update offset data: %s", err.Error())
		}
	}
}

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
