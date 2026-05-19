# Ingestion Node

The ingestion node is a centralized service responsible for collecting logs from multiple logging agents, processing and enriching them, and forwarding them for downstream handling.

---

## How it works

Logging agents send log data to the ingestion node via the `/ingest` endpoint. The node first validates that the incoming logs are properly formatted. If the logs are invalid, it returns an error response.

Once validated, the logs are enriched with the sender’s IP address and then persisted by writing them to a log file for storage and later processing.

---

## How to use

1. Configure the `config.yaml` file to specify the destination file path for storing logs.
2. Update the `logging agent’s config.yaml` to call the the ingestion node’s `/ingest` endpoint.
3. Start the ingestion service:

```bash
go run cmd/main.go
```

## Features

- Centralized log collection
- Log validation 
- Log enrichment
    - Automatically adds metadata to logs, such as the sender’s IP address.

- Persistent storage  
  - Writes processed logs to a configured log file for durability and later analysis.

- Configurable setup
  - Supports external configuration via `config.yaml`, including log file path and other runtime settings.

- Easy Agent integration  
  - Easily integrates with logging agents through endpoint configuration.

## Current Design Choices

- **Minimal log enrichment (IP only)**  
  The logger currently enriches logs only with the sender’s IP address. This design decision was made because the system was built primarily for learning and local development purposes. Since it runs in a localhost environment, there is limited contextual information available for deeper enrichment (such as geo-location, environment metadata, or service topology).

- **Local log persistence**  
  Logs are persisted locally to a file rather than being sent to external storage systems or message queues. This approach was chosen because the expected log volume is low, and the system is intended for simplicity and educational use. Introducing message queues or distributed storage at this stage would add unnecessary complexity without providing proportional benefit.

## Future Scope of Development

- **Extended configuration options**

  Enhance the node’s configuration capabilities to support multiple log delivery modes, including:
  - Local file-based persistence
  - External persistent storage systems (e.g., databases or object storage)
  - Message queues or streaming platforms for downstream processing

- **At-least-once delivery guarantee**

  Improve reliability by implementing at-least-once delivery semantics for forwarding logs to downstream systems. This ensures logs are not lost in transit, even in the event of failures or retries.

- **Enhanced log enrichment**

  Expand the enrichment pipeline to include additional contextual and diagnostic metadata, such as:
  - Request headers and user-agent information  
  - Service or application identifiers  
  - Hostname and environment (e.g., staging, production)  
  - Request latency and processing time  
  - Geolocation data based on IP address (if applicable)  
  - Correlation or trace IDs for distributed tracing  