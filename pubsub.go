package audiy_importer


import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

// pubsubClient is pubsub client wrapper
type pubsubClient struct {
	topic  string
	client *pubsub.Client
}

func NewClient(project, topic string) (*pubsubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	psCli := pubsubClient{
		topic:  topic,
		client: client,
	}
	return &psCli, err
}

func (c *pubsubClient) Publish(ctx context.Context, serialized []byte) error {
	topic := c.client.Topic(c.topic)
	result, err := topic.Publish(ctx, &pubsub.Message{Data: serialized}).Get(ctx)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	log.Printf("successfully published ID: %v", result)
	return nil
}

//
//
//func configurePubsub(projectID string) (*pubsub.Client, error) {
//
//	PubsubTopicID := os.Getenv("PUBSUB_TOPIC")
//	ctx := context.Background()
//	client, err := pubsub.NewClient(ctx, projectID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Create the topic if it doesn't exist.
//	if exists, err := client.Topic(PubsubTopicID).Exists(ctx); err != nil {
//		return nil, err
//	} else if !exists {
//		if _, err := client.CreateTopic(ctx, PubsubTopicID); err != nil {
//			return nil, err
//		}
//	}
//	return client, nil
//}


