# REST API - Receive messages from LINE and Facebook Messenger

## Overview
Receive messages from LINE and Facebook Messenger, validate the message format/schema, and send the message to RabbitMQ queue.

```mermaid
sequenceDiagram
    participant LINE as LINE
    participant API as REST API
    participant Caching as Redis Caching
    participant RMQ as RabbitMQ
    LINE->>API: Send message
    API->>Caching: Check if message is duplicate
    Caching-->>API: Return duplicate status
    alt If message is duplicate
        API->>LINE: Return duplicate message response
    else If message is not duplicate
        API->>RMQ: Send message to RabbitMQ queue
        RMQ-->>API: Return message sent status
        API->>LINE: Return success response
    end
``` 

## Workflow
1. Recieve message with LINE to API endpoint
2. Convert message to LINE message object
3. Validate message format/schema from system schema file
4. Send message to RabbitMQ queue
   - Queue name: `line-messages`

## Observability with OpenTelemetry
- Trace the entire workflow from receiving the message to sending it to RabbitMQ queue
- Log the message content and processing status at each step
- Monitor the performance of the API endpoint and RabbitMQ queue

## Input of LINE message object
* LINE message object

### LINE message object

Text Message Format
```json
{
  "to": "USER_ID",
  "messages": [
    {
      "type": "text",
      "text": "Hello, world!"
    }
  ]
}
```

Image Message Format
```json
{
  "to": "USER_ID",
  "messages": [
    {
      "type": "image",
      "originalContentUrl": "https://example.com",
      "previewImageUrl": "https://example.com"
    }
  ]
}
```

## RabbitMQ with Work Queues pattern 
- Queue name: `line-messages`

### Message format in RabbitMQ queue
```json
{
  "name": "line-messages",
  "message": {
    "to": "USER_ID",
    "messages": [
      {
        "type": "text",
        "text": "Hello, world!"
      }
    ]
  }
}
```

