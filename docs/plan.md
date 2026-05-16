# Distributed Logger

## components
### 1. Logger SDK
- Integrated into application services
- Generates structured logs
- Writes logs locally (file/stdout)
- Independent from distributed logging infrastructure

### 2. Logging Agent
- Monitors log sources (files/stdout)
- Buffers and batches logs
- Retries failed transmissions
- Sends logs to ingestion service

### 3. Ingestion Service
- Receives logs from agents
- Validates and forwards logs to the message broker

### 4. Message Queue / Broker
- Buffers logs
- Decouples producer and consumer throughput
- Handles burst traffic

### 5. Logging Processor Service
- Consumes logs from queue
- Processes and enriches logs
- Stores logs persistently

### 6. Persistent Storage
- Stores logs durably
- Supports querying and retention policies

## phases

### phase 1
- Implement the logging service for handling logs
    - The logging service should have a single client connection
    - Store logs in an AOF periodically
- Implement the logger sdk
    - Implement the logger to have base implementation plus allow services to extend it
- Implement the logging agent to pickup and transmit emitted logs to the logging service
    - should monitor a configured file and send logs to the logging service
- Integrate the sdk into a sample service
    - Integrate the logger sdk into a service
    - Verify the logging in the logging service and the AOF file
    - Verify delivery guarantees in the logs

### phase 2
- Make the logging agent send data to the ingestion node instead of the logging service
- Implement the ingestion node
    - This will get data from the logging agents
    - Verify the payload and send to message queue for downstream processing
- Implement the message queue
    - The message queue takes in data from the ingestion nodes and the logging service will pull the data from here
    - Helpful in decoupling producer speed and consumer speed
- Scale the logging service
    - Make the logging service diconnected from the client and fetch data only from the message queue
    - Upgrade the AOF implementation to WAL for more durability
- Integrate the sdk into multiple sample services
    - Integrate the logger sdk into multiple services
- Deliver ordering guarantees
    - Verify the delivery guarantee and ordering in the WAL
        - per client ordering guarantees

### phase 3
- To be done

## design choices
- Choice of persistent storage
    - `mongodb`  
        - Advantages  
            - Easier to manage
            - Need not worry about data corruption and storage
            - Easier to scale  
        - Disadvantage
            - Would not understand how to handle corruption, log sync, etc
    - custom `AOF`
        - Advantages
            - Really good learning opportunity
            - Will understand in depth how to handle corruption, sync issues
        - Disadvantages
            - Very hard to implement and test for correctness
- Choice of language
    - `go`
        - Advantages
            - Natural Concurrency
            - Easier to start with for distributed systems
        - Disadvantages
            - Not really much, good for distributed systems

## Testing
- should test logging service failure scenarios
    - verify replay correctness on restart
    - verify acknowledged logs are not lost
- WAL checks
    - verify partial writes do not corrupt WAL
    - verify partial/corrupted writes do not lead to incorrect logs
