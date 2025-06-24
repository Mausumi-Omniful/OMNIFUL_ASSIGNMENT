package pusher

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/omniful/go_commons/log"
	beams "github.com/pusher/push-notifications-go"
)

type Publisher interface {
	PublishToWeb(request *PublishRequest) (string, error)
}

type BeamsClient struct {
	beams.PushNotifications
}

const (
	interestWebKey = "web"

	interestDelimiter = "-"
)

// PublishRequest is the data that the system accepts at client side.
// WebConfig for push notifications
type PublishRequest struct {
	Data      WebConfig `json:"data,omitempty"`
	Interests []string  `json:"interests" validate:"required"`
}

// WebConfig
// Notification holds details specific to the appearance of the web push notification
// Dir is the text direction for the push notification (e.g., "ltr" or "rtl")
// Tag is a tag for the push notification, useful for grouping related notifications
// Vibrate specifies the vibration pattern for the device when the notification is received
type WebConfig struct {
	Notification WebNotification `json:"notification,omitempty"`
	Title        string          `json:"title,omitempty"`
	Body         string          `json:"body,omitempty"`
	Icon         string          `json:"icon,omitempty"`
	Badge        string          `json:"badge,omitempty"`
	Image        string          `json:"image,omitempty"`
	Dir          string          `json:"dir,omitempty"`
	Lang         string          `json:"lang,omitempty"`
	Tag          string          `json:"tag,omitempty"`
	Timestamp    int64           `json:"timestamp,omitempty"`
	Vibrate      []any           `json:"vibrate,omitempty"`
}

type WebNotification struct {
	Title       string                `json:"title,omitempty"`
	Body        string                `json:"body,omitempty"`
	Messages    []NotificationMessage `json:"messages,omitempty"`
	Meta        map[string]any        `json:"meta"`
	IsMetaExist bool                  `json:"is_meta_exist"`
}

type NotificationMessage struct {
	NotificationType ActionType     `json:"notification_type"`
	MessageType      ContentType    `json:"message_type"`
	Content          string         `json:"content"`
	Config           map[string]any `json:"config"`
	Sequence         int            `json:"sequence"`
}

func (c *BeamsClient) PublishToWeb(request *PublishRequest) (string, error) {
	if err := c.validateRequest(request); err != nil {
		log.Errorf("[PublishToWeb] Error while validating request :: %+v | err :: %s ", request, err.Error())
		return "", err
	}

	publishId, err := c.PublishToInterests(request.Interests, map[string]any{interestWebKey: request})
	if err != nil {
		log.Errorf(
			"[PublishToWeb] Error while publishing to interests request :: %+v | err :: %s",
			request,
			err.Error(),
		)
		return "", err
	}

	return publishId, nil
}

func (c *BeamsClient) validateRequest(request *PublishRequest) error {
	if err := validator.New().Struct(request); err != nil {
		return errors.New("[validateRequest] request invalid")
	}

	if len(request.Interests) == 0 {
		return errors.New("[validateRequest] interests is empty")
	}

	return nil
}

func (p *PublishRequest) NewWebNotification(
	interest string,
	tenantID string,
	web WebConfig,
) *PublishRequest {
	return &PublishRequest{
		Interests: []string{p.createInterest(interest, tenantID)},
		Data:      web,
	}
}

func (p *PublishRequest) createInterest(interestName, tenantID string) string {
	return interestName + interestDelimiter + tenantID
}

func NewPublisherRequest() *PublishRequest {
	return &PublishRequest{}
}

func NewBeamsClient(instanceId, secretKey string) *BeamsClient {
	beamsClient, err := beams.New(instanceId, secretKey)
	if err != nil {
		log.Errorf("[NewBeamsClient] Error while creating beams client :: %s", err.Error())
		return nil
	}

	return &BeamsClient{
		PushNotifications: beamsClient,
	}
}
