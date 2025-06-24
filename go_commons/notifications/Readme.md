# Notifications Package

## Overview
The notifications package provides a flexible and extensible framework for sending notifications across various channels. Currently, it implements Slack notifications with rich formatting and interactive capabilities, with the architecture supporting easy addition of other notification channels.

## Features
- Interface-based design for multiple notification channels
- Comprehensive Slack integration supporting:
  - Basic text messages
  - Thread replies
  - Rich message attachments
  - Interactive buttons
  - Message updates

## Installation
```go
import "github.com/omniful/go_commons/notifications"
```

## Components

### Core Interface
```go
type Notifications interface {
    SendNotification(message, channelID string) (err error)
}
```

### Slack Implementation
The Slack implementation provides rich functionality for sending Slack messages:

```go
package main

import (
    "github.com/omniful/go_commons/notifications/slack"
)

func main() {
    // Initialize Slack client
    slackClient := slack.NewClient("your-slack-token")
    
    // Send basic notification
    err := slackClient.SendNotification("Hello World!", "channel-id")
    if err != nil {
        // Handle error
    }
    
    // Send message with interactive buttons
    attachment := slack.NewAttachment(
        "Would you like to proceed?",
        slack.PlainText,
        []slack.Actions{
            {
                ID:          "btn_yes",
                Value:       "yes",
                Text:        "Yes",
                ElementType: slack.PlainText,
                Style:       "primary",
            },
            {
                ID:          "btn_no",
                Value:       "no",
                Text:        "No",
                ElementType: slack.PlainText,
                Style:       "danger",
            },
        },
    )
    
    err = slackClient.SendNotificationWithAttachment(
        "Action Required",
        "channel-id",
        attachment,
    )
}
```

## Advanced Usage

### Thread Replies
```go
// Send a reply in a thread
err := slackClient.SendNotificationInThread(
    "This is a reply",
    "parent-message-timestamp",
    "channel-id",
)
```

### Updating Messages
```go
// Update an existing message
err := slackClient.UpdateNotificationWithAttachment(
    "Updated message",
    "channel-id",
    "message-timestamp",
    updatedAttachment,
)
```

## Message Formatting
The package supports two types of text formatting:
- `PlainText`: Regular text without markdown
- `MarkDown`: Markdown-formatted text

## Error Handling
All methods return an error type that should be checked for successful operation. The package includes built-in logging for error cases.

## Best Practices
1. Always check for errors returned by notification methods
2. Store Slack tokens securely (e.g., environment variables)
3. Use appropriate message formatting based on content type
4. Consider rate limits when sending multiple notifications

## Notes
- The package is designed to be thread-safe
- Future implementations may include other notification channels (email, SMS, etc.)
- The Slack implementation uses the official `slack-go/slack` package
