package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"google.golang.org/api/option"
)

var Module = fx.Options(
	fx.Provide(NewClient),
)

type Client struct {
	googleClient *pubsub.Client
}

func NewClient(c config.Config) (*Client, error) {
	googleClient, err := pubsub.NewClient(context.Background(), c.GCP.ProjectID, option.WithoutAuthentication())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pubsub client")
	}

	return &Client{
		googleClient: googleClient,
	}, nil

}

func (c *Client) Publish(ctx context.Context, topic string, msg []byte) error {

	t, err := getOrCreateTopic(ctx, c, topic)
	if err != nil {
		return errors.Wrap(err, "failed to get or create topic")
	}

	res := t.Publish(ctx, &pubsub.Message{Data: msg})

	if _, err := res.Get(ctx); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, topic string, subscriberID string, handler func(ctx context.Context, msg *pubsub.Message)) error {

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

func getOrCreateTopic(ctx context.Context, client *Client, topicID string) (*pubsub.Topic, error) {
	t := client.googleClient.Topic(topicID)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if topic exists")
	}
	if ok {
		return t, nil
	}
	return client.googleClient.CreateTopic(ctx, topicID)
}

func getOrCreateSubscription(ctx context.Context, client *Client, topicID string, subscriberID string) (*pubsub.Subscription, error) {
	t, err := getOrCreateTopic(ctx, client, topicID)
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
