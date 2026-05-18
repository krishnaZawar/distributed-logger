# Logging Agent

The Logging Agent is a background service that collects logs from configured log files and forwards them to a remote endpoint for storage and further processing.

It is designed to work with `.ndjson` (newline-delimited JSON) log files.

---

## How It Works

The agent continuously monitors configured log files, detects newly appended log entries, processes them in batches, and reliably delivers them to a remote endpoint.

It is designed to run indefinitely as a background service.

---

## How to Use

### 1. Configure `config.yaml`

Before starting the agent, ensure `config.yaml` is properly set up:

- Specify the log files to monitor
  - Ensure the files already exist at the provided paths
  - Ensure the files are in `.ndjson` format (invalid formats may lead to unexpected behavior)
- Specify the path for the agent’s internal log file
- Configure the remote endpoint details for log delivery
- Set the batch size for processing logs per file

For more information on this, read the config.yaml file

**Note: All configuration values must be set before starting the agent.**


---

### 2. Start the Logging Agent

```bash
go run cmd/main.go
```
Now the agent is up and running indefinitely

### 3. For addition of new log files for tracking

To add new files for tracking:
- Update the config.yaml with the new log files
- Restart the log agent

**Note: This wont affect the log tracking for older files and agent will resume from where it left**

---

## Features
- At-least-once delivery guarantee for log ingestion
- Supports tracking multiple log files concurrently
- Resumes from the last processed offset after restart
- Thread-safe processing for concurrent file handling
- Batch-based log processing for improved efficiency

---

## Failure Handling

The Logging Agent is designed to be resilient to common system and network failures. It prioritizes reliability and data preservation over strict delivery uniqueness.

### Endpoint is Down

When the remote logging endpoint is unavailable:

- The agent retries delivery after a timeout
- If the logs are not delivered even in the retry, those logs are processed again in the next iteration
- Logs are **not discarded**

### Agent Crash or Restart

When the Logging Agent crashes or is restarted:

- The last processed log position (offset) is retrieved from persistent storage
- The agent resumes reading from the last saved offset
- This ensures logs are not skipped after restart

**Note: Some logs may be reprocessed, which can lead to duplicates.**

---

## Guarantees

The Logging Agent provides the following delivery and processing guarantees:

### At-Least-Once Delivery

- Every log is guaranteed to be delivered **at least once**
- In failure scenarios, the same log may be delivered multiple times
- This ensures no data loss even under crashes or network failures

---

### No Guaranteed Ordering Across Multiple Files

- Ordering is preserved **only within a single file**
- There is no global ordering guarantee across multiple log files

---

## Future Scope of Development
1. Extensible support for multiple log formats beyond ".ndjson" files
2. Dynamic configuration reload without restarting the agent