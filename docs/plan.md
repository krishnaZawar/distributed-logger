# Distributed Logger

## components
- logger sdk
    - integrated service side, generates and transmits the logs
- ingestion node (worker queue)
    - collects the logs from the services and stores them for use by logging service
    - essentially used for decoupling producer speed and consumer speed
- logging service
    - collects logs from the ingestion node and store in a persistent store
- persistent storage
    - holds the logs of services implementing the logger sdk

## phases

### phase 1
- Implement the logging service for handling logs
    - The logging service should have a single client connection
    - Store logs in an AOF periodically
- Implement the logger sdk
    - Implement the logger to have base implementation plus allow services to extend it
    - Make the logger resilient to disconencts, i.e., disconnects should lead to retries and the pending logs should be buffered for that moment
- Integrate the sdk into a sample service
    - Integrate the logger sdk into a service
    - Verify the logging in the logging service and the AOF file
    - Verify delivery guarantees in the logs

### phase 2
- Scale the logging service
    - Make the logging service diconnected from the client and implement an ingestion node in the middle for log handling
    - Deliver ordering guarantees
    - Upgrade the AOF implementation to WAL for more durability
- Integrate the sdk into multiple sample services
    - Integrate the logger sdk into multiple services
    - Verify the delivery guarantee and ordering in the WAL
        - per client ordering guarantees
        - other too exist, need to decide on the needed semantic

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
