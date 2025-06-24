# Package pusher

## Overview
The pusher package provides a robust real-time notification system with multiple implementations for different use cases. It supports both Pusher Channels and Pusher Beams for web push notifications, offering flexible ways to send real-time updates to clients.

## Features
- **Pusher Channels**: Real-time messaging through channels
- **Pusher Beams**: Web push notifications
- **Multiple Event Types**: Support for various notification types (Toast, Download, Navigate, etc.)
- **Flexible Message Formatting**: Structured message format with metadata support
- **Channel-based Communication**: Target specific users or groups through channels

## Components

### 1. Pusher Client (V1)
Basic implementation for channel-based messaging.

```go
package main

import (
	"github.com/omniful/go_commons/pusher"
)

func main() {
	// Initialize client
	client := pusher.NewClient("app-id", "key", "secret", "cluster")
	
	// Send a toast notification
	event := pusher.NewToastEvent(
		"notification_event",
		"Hello World!",
		pusher.SuccessMessage,
	)
	
	// Send to specific channel and user
	err := client.SendEventMessage("channel-name", "user-id", event)
	if err != nil {
		// Handle error
	}
	
	// Send download notification
	downloadEvent := pusher.NewDownloadEvent(
		"download_event",
		"Your file is ready",
		"https://download-url.com",
		pusher.InfoMessage,
	)
	
	err = client.SendEventMessage("channel-name", "user-id", downloadEvent)
	if err != nil {
		// Handle error
	}
}
```

### 2. Pusher Client V2
Enhanced implementation with more structured events and actions.

```go
package main

import (
	"github.com/omniful/go_commons/pusher"
)

func main() {
	client := pusher.NewClient("app-id", "key", "secret", "cluster")
	
	// Create actions
	actions := []pusher.Action{
		{
			ActionType: pusher.ToastActionType,
			Data: pusher.Data{
				Type: pusher.MessageDataType,
				Message: pusher.Message{
					Content:     "Operation completed successfully",
					ContentType: pusher.SuccessContentType,
				},
			},
			Sequence: 1,
		},
		{
			ActionType: pusher.NavigateActionType,
			Data: pusher.Data{
				Type: pusher.PathDataType,
				Path: "/dashboard",
			},
			Sequence: 2,
		},
	}
	
	// Create notification event
	event := pusher.NewNotificationEvent(
		"custom_event",
		pusher.CustomEvent,
		true,
		actions,
	)
	
	// Send event
	err := client.SendPusherEvent("channel-name", "user-id", event)
	if err != nil {
		// Handle error
	}
}
```

### 3. Pusher Beams
Web push notifications implementation.

```go
package main

import (
	"github.com/omniful/go_commons/pusher"
)

func main() {
	// Initialize Beams client
	beamsClient := pusher.NewBeamsClient("instance-id", "secret-key")
	
	// Create web notification
	webConfig := pusher.WebConfig{
		Title: "New Notification",
		Body:  "You have a new message",
		Icon:  "https://example.com/icon.png",
		Notification: pusher.WebNotification{
			Title: "New Message",
			Body:  "Hello from Pusher Beams!",
			Messages: []pusher.NotificationMessage{
				{
					NotificationType: pusher.InfoActionType,
					MessageType:      pusher.NotificationContentType,
					Content:         "Your order has been processed",
					Sequence:        1,
				},
			},
		},
	}
	
	// Create publish request
	req := pusher.NewPublisherRequest()
	notification := req.NewWebNotification(
		"notification",
		"tenant-123",
		webConfig,
	)
	
	// Publish notification
	publishID, err := beamsClient.PublishToWeb(notification)
	if err != nil {
		// Handle error
	}
}
```

## Event Types
- `TOAST`: Simple toast messages
- `DOWNLOAD_LINK`: File download notifications
- `NAVIGATE`: Navigation events
- `INFO`: Information messages
- `AUDIO_NOTIFICATION`: Audio notifications
- `DATA`: Custom data events

## Message Types
- `success`: Success messages
- `warning`: Warning messages
- `failure`: Error messages
- `info`: Information messages
- `in_progress`: Progress updates
- `notification`: General notifications

## Best Practices
1. Always handle errors returned by the client methods
2. Use appropriate event types based on the notification purpose
3. Set proper sequence numbers when using multiple actions
4. Include relevant metadata when needed
5. Use channel names that follow a consistent pattern
6. Set `IsLastResponse` appropriately to manage client-side subscriptions

## Notes
- The package requires valid Pusher credentials (App ID, Key, Secret, Cluster)
- For Beams implementation, valid Instance ID and Secret Key are required
- Channel names are automatically formatted as `{channel}-{ID}`
- Events can be targeted to specific users or groups using the ID parameter
