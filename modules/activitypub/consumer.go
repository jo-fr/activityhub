package activitypub

import (
	"context"
	"encoding/json"

	"github.com/jo-fr/activityhub/pkg/log"

	"cloud.google.com/go/pubsub"

	"github.com/jo-fr/activityhub/modules/feed/model"
	"github.com/jo-fr/activityhub/pkg/externalmodel"
	te "github.com/jo-fr/activityhub/pkg/pubsub"
)

type Consumer struct {
	handler *Handler
	pubsub  *te.Client
	log     *log.Logger
}

func NewConsumer(log *log.Logger, pubsub *te.Client, handler *Handler) *Consumer {
	return &Consumer{
		handler: handler,
		pubsub:  pubsub,
		log:     log,
	}
}

func (c *Consumer) consumeOutbox() func(ctx context.Context, msg *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		var status model.Status
		if err := json.Unmarshal(msg.Data, &status); err != nil {
			c.log.Errorf("error while unmarshalling message: %s", err.Error())
			msg.Nack()
			return
		}

		if err := c.handler.SendPostToFollowers(ctx, status.AccountID, status.Content); err != nil {
			c.log.Errorf("error while sending post to followers: %s", err.Error())
			msg.Nack()
			return
		}

		msg.Ack()
	}
}

func (c *Consumer) consumeInbox() func(ctx context.Context, msg *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		var activity externalmodel.Activity
		if err := json.Unmarshal(msg.Data, &activity); err != nil {
			c.log.Errorf("error while unmarshalling message: %s", err.Error())
			msg.Nack()
			return
		}

		err := c.handler.ReceiveInboxActivity(ctx, activity)
		if err != nil {
			c.log.Errorf("error while receiving inbox activity: %s", err.Error())
			msg.Nack()
			return
		}

		msg.Ack()
	}
}

func Subscribe(consumer *Consumer, handler *Handler, logger *log.Logger) error {
	go func() {
		err := consumer.pubsub.Subscribe(context.Background(), te.TopicOutbox, "activitpub_module_outbox_consumer", consumer.consumeOutbox())
		if err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		err := consumer.pubsub.Subscribe(context.Background(), te.TopicInbox, "activitpub_module_inbox_consumer", consumer.consumeInbox())
		if err != nil {
			logger.Error(err)
		}
	}()

	return nil
}
