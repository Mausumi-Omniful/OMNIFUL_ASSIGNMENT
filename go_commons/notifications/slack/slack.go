package slack

import (
	"bytes"

	"github.com/omniful/go_commons/log"
	"github.com/slack-go/slack"
)

type Slack struct {
	slack *slack.Client
}

type ElementType string

const (
	PlainText ElementType = "plain_text"
	MarkDown  ElementType = "mrkdwn"
)

type Attachment struct {
	Text        string
	Actions     []Actions
	ElementType ElementType
}

type Actions struct {
	ID          string
	Value       string
	Text        string
	ElementType ElementType
	Style       slack.Style
}

func NewClient(token string) *Slack {
	if token == "" {
		panic("not able to find access token for slack")
	}

	return &Slack{slack: slack.New(token)}
}

func NewAttachment(text string, elementType ElementType, actions []Actions) Attachment {
	return Attachment{
		Text:        text,
		ElementType: elementType,
		Actions:     actions,
	}
}

func (s *Slack) SendNotification(message, channelID string) (err error) {
	_, _, err = s.slack.PostMessage(channelID, slack.MsgOptionText(message, true))
	if err != nil {
		log.Errorf("unable to send slack notification for this message : %s and channelID : %s. err :: %v",
			message, channelID, err.Error(),
		)
		return
	}

	return
}

func (s *Slack) SendNotificationInThread(message, timestamp, channelID string) (err error) {
	_, _, err = s.slack.PostMessage(channelID, getMessageParams(message, nil, timestamp)...)
	if err != nil {
		log.Errorf("unable to send slack notification for this message in thread: %s and channelID : %s. err :: %v",
			message, channelID, err.Error(),
		)
		return
	}

	return
}

func (s *Slack) SendNotificationWithAttachment(message, channelID string, request Attachment) (err error) {
	_, _, err = s.slack.PostMessage(channelID, getMessageParams(message, &request, "")...)
	if err != nil {
		log.Errorf("unable to send slack notification for this message : %s and channelID : %s. err :: %v",
			message, channelID, err.Error(),
		)
		return
	}

	return
}

func (s *Slack) UpdateNotificationWithAttachment(
	message,
	channelID string,
	timestamp string,
	request Attachment,
) (err error) {
	_, _, _, err = s.slack.UpdateMessage(channelID, timestamp, getMessageParams(message, &request, "")...)
	if err != nil {
		log.Errorf("unable to update slack notification for this message : %s and channelID : %s. err :: %v",
			message, channelID, err.Error(),
		)
		return
	}

	return
}

func (attachment Attachment) getBlockElement() (blockElement []slack.BlockElement) {
	for _, action := range attachment.Actions {
		blockElement = append(blockElement, slack.NewButtonBlockElement(
			action.ID,
			action.Value,
			slack.NewTextBlockObject(string(action.ElementType), action.Text, false, false),
		).WithStyle(action.Style))
	}
	return
}

func (attachment Attachment) getSlackAttachment() slack.Attachment {
	sectionBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(string(attachment.ElementType), attachment.Text, false, false),
		nil,
		nil,
	)

	slackAttachment := slack.Attachment{
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{sectionBlock},
		},
	}

	if len(attachment.Actions) > 0 {
		buttonBlock := slack.NewActionBlock("", attachment.getBlockElement()...)
		slackAttachment.Blocks.BlockSet = append(slackAttachment.Blocks.BlockSet, buttonBlock)
	}

	return slackAttachment
}

func getMessageParams(message string, attachment *Attachment, threadTimestamp string) []slack.MsgOption {
	var options []slack.MsgOption

	if len(threadTimestamp) > 0 {
		options = append(options, slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			AsUser:          true,
			ThreadTimestamp: threadTimestamp,
		}))
	} else {
		options = append(options, slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			AsUser: true,
		}))
	}

	if len(message) > 0 {
		options = append(options, slack.MsgOptionText(message, false))
	}

	if attachment != nil {
		options = append(options, slack.MsgOptionAttachments(attachment.getSlackAttachment()))
	}

	return options
}

func (s *Slack) SendNotificationWithFile(
	message,
	channelID,
	filename string,
	fileContent []byte,
) (err error) {
	_, err = s.slack.UploadFileV2(slack.UploadFileV2Parameters{
		Filename:       filename,
		Reader:         bytes.NewReader(fileContent),
		Channel:        channelID,
		InitialComment: message,
		FileSize:       len(fileContent),
	})
	if err != nil {
		return
	}

	return
}
