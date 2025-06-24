package pusher

import (
	"encoding/json"
)

type (
	PusherEvent string
	ActionType  string
	DataType    string
	ContentType string
)

const (
	CustomEvent PusherEvent = "CUSTOM"
)

const (
	InfoActionType        ActionType = "INFO"
	ToastActionType       ActionType = "TOAST"
	DownloadActionType    ActionType = "DOWNLOAD"
	NavigateActionType    ActionType = "NAVIGATE"
	AudioNotificationType ActionType = "AUDIO_NOTIFICATION"
	DataActionType        ActionType = "DATA"
)

const (
	MessageDataType DataType = "message"
	URLDataType     DataType = "url"
	PathDataType    DataType = "path"
)

const (
	InfoContentType         ContentType = "info"
	SuccessContentType      ContentType = "success"
	WarningContentType      ContentType = "warning"
	FailureContentType      ContentType = "failure"
	InProgressContentType   ContentType = "in_progress"
	NotificationContentType ContentType = "notification"
)

// NotificationEvent is the data that the system accepts at client side.
// IsLastResponse signify in case it is false FE will unsubscribe this event
// **PusherEvent is always CUSTOM for now**
type NotificationEvent struct {
	Name           string      `json:"event_name" validate:"required"`
	Type           PusherEvent `json:"event_type" validate:"required"`
	IsLastResponse bool        `json:"is_last_response"`
	Actions        []Action    `json:"actions"`
}

type ActionList []Action

// Action Type is the type of action that the system will take
// e.g. INFO, TOAST, DOWNLOAD, NAVIGATE
// Sequence is the order of the action list that FE will use to execute the actions
type Action struct {
	ActionType ActionType `json:"action_type"`
	Data       Data       `json:"data"`
	Sequence   int        `json:"sequence"`
}

// Data Type represent the information type that the system will send to FE
// e.g. message, url, path
// If DataType is 'message' then go to 'message' field
// If DataType is 'url' then go to 'url' field
// If DataType is 'path' then go to 'path' field
type Data struct {
	Type    DataType `json:"type"`
	Url     string   `json:"url"`
	Path    string   `json:"path"`
	Message Message  `json:"message"`
}

type Message struct {
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
	Meta        Meta        `json:"meta"`
	IsMetaExist bool        `json:"is_meta_exist"`
}

// Meta can be used to send any additional data to FE
type Meta map[string]any

func NewNotificationEvent(
	eventName string,
	eventType PusherEvent,
	isLastResponse bool,
	actions ActionList,
) *NotificationEvent {
	return &NotificationEvent{
		Name:           eventName,
		Type:           eventType,
		IsLastResponse: isLastResponse,
		Actions:        actions,
	}
}

func (c *Client) SendPusherEvent(channel, ID string, event *NotificationEvent) error {
	data, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		return marshalErr
	}

	return c.PushMessageToChannel(channel, event.Name, ID, data)
}
