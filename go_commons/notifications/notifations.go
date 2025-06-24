package notifications

type Notifications interface {
	SendNotification(message, channelID string) (err error)
}
