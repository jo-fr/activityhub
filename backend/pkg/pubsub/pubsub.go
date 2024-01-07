package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewClient),
)

type Client struct {
	googleClient *pubsub.Client
}

func NewClient(lc fx.Lifecycle, c config.Config, log *log.Logger) (*Client, error) {

	googleClient, err := pubsub.NewClient(context.Background(), c.GCP.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pubsub client")
	}

	client := &Client{
		googleClient: googleClient,
	}

	registerHooks(lc, client, log)
	log.Info("pubsub client created")
	return client, nil

}

// registerHooks for uber fx
func registerHooks(lc fx.Lifecycle, c *Client, logger *log.Logger) {
	lc.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				logger.Info("closing pubsub connection")
				return c.googleClient.Close()
			},
		},
	)
}

func (c *Client) Publish(ctx context.Context, topic Topic, msg any) error {

	json, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	t, err := getOrCreateTopic(ctx, c, topic)
	if err != nil {
		return errors.Wrap(err, "failed to get or create topic")
	}

	res := t.Publish(ctx, &pubsub.Message{Data: json})

	if _, err := res.Get(ctx); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, topic Topic, subscriberID string, handler func(ctx context.Context, msg *pubsub.Message)) error {

	sub, err := getOrCreateSubscription(ctx, c, topic, subscriberID)
	if err != nil {
		return errors.Wrap(err, "failed to get or create subscription")
	}

	err = sub.Receive(ctx, handler)
	if err != nil {
		return errors.Wrap(err, "failed to receive message")
	}

	return nil
}

func getOrCreateTopic(ctx context.Context, client *Client, topicID Topic) (*pubsub.Topic, error) {
	t := client.googleClient.Topic(topicID.String())
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if topic exists")
	}
	if ok {
		return t, nil
	}
	return client.googleClient.CreateTopic(ctx, topicID.String())
}

func getOrCreateSubscription(ctx context.Context, client *Client, topicID Topic, subscriberID string) (*pubsub.Subscription, error) {
	t, err := getOrCreateTopic(ctx, client, Topic(topicID.String()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get or create topic")
	}

	s := client.googleClient.Subscription(subscriberID)
	ok, err := s.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if subscription exists")
	}
	if ok {
		return s, nil
	}
	return client.googleClient.CreateSubscription(ctx, subscriberID, pubsub.SubscriptionConfig{Topic: t})
}
