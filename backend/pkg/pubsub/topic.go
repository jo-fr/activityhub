package pubsub

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	TopicInbox  Topic = "activityhub_inbox_queue"
	TopicOutbox Topic = "activityhub_outbox_queue"
)
