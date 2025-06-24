package pusher

import (
	"encoding/json"
	"errors"
	"github.com/omniful/go_commons/log"
	"github.com/pusher/pusher-http-go"
)

type Client struct {
	pusher.Client
}

type EventType string
type EventMessageType string

const (
	ToastEventType        EventType = "TOAST"
	DownloadFileEventType EventType = "DOWNLOAD_LINK"

	SuccessMessage EventMessageType = "success"
	WarningMessage EventMessageType = "warning"
	FailureMessage EventMessageType = "failure"
	InfoMessage    EventMessageType = "info"
)

/*
Event is the data that the system accepts at client side.
IsLastResponse signify in case the system need to send the progress event to FE for example to show export job process
and when completed will send the last event with is_last_response=true
*/
type Event struct {
	Name           string    `json:"event_name"`
	Type           EventType `json:"event_type"`
	Data           EventData `json:"data"`
	IsLastResponse bool      `json:"is_last_response"`
}

type EventData struct {
	Message *EventMessage `json:"message"`
	Link    string        `json:"link"`
}

type EventMessage struct {
	Type    EventMessageType `json:"type"`
	Content string           `json:"content"`
	Info    interface{}      `json:"info"`
}

func NewToastEvent(eventName, displayMessage string, messageType EventMessageType, info ...interface{}) *Event {
	return &Event{
		Name:           eventName,
		Type:           ToastEventType,
		IsLastResponse: true,
		Data: EventData{
			Message: &EventMessage{
				Type:    messageType,
				Content: displayMessage,
				Info:    info,
			},
		},
	}
}

func NewDownloadEvent(eventName, displayMessage, link string, messageType EventMessageType, info ...interface{}) *Event {
	return &Event{
		Name:           eventName,
		Type:           DownloadFileEventType,
		IsLastResponse: true,
		Data: EventData{
			Message: &EventMessage{
				Type:    messageType,
				Content: displayMessage,
				Info:    info,
			},
			Link: link,
		},
	}
}

func NewClient(appID, key, secret, cluster string) *Client {
	c := pusher.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  true,
	}

	return &Client{
		Client: c,
	}
}

// PushMessageToChannel : Data should be json marshal value
// ID : Use the user ID to push the message to a specific user, or the tenant ID to push the message to multiple users within a tenant.
func (c *Client) PushMessageToChannel(channel, eventName, ID string, data interface{}) error {
	if ID == "" {
		return errors.New("ID cannot be empty")
	}

	channelName := c.channelName(channel, ID)
	err := c.Trigger(channelName, eventName, data)
	if err != nil {
		log.Errorf(
			"Error in pushing message to channel: %s, eventName: %s, data: %+v with error: %s",
			channelName,
			eventName,
			data,
			err.Error(),
		)
		return err
	}

	return nil
}

func (c *Client) SendEventMessage(channel, ID string, event *Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return c.PushMessageToChannel(channel, event.Name, ID, data)
}

func (c *Client) channelName(channel, ID string) string {
	return channel + "-" + ID
}
