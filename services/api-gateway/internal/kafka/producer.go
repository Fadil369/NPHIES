package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Producer handles Kafka message publishing
type Producer struct {
	writers map[string]*kafka.Writer
	logger  *logrus.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, logger *logrus.Logger) (*Producer, error) {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		logger:  logger,
	}, nil
}

// getWriter gets or creates a Kafka writer for a specific topic
func (p *Producer) getWriter(topic string, brokers []string) *kafka.Writer {
	if writer, exists := p.writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
		ErrorLogger:            kafka.LoggerFunc(p.logger.Errorf),
		Compression:            kafka.Snappy,
		AllowAutoTopicCreation: true,
	}

	p.writers[topic] = writer
	return writer
}

// Publish publishes a message to a Kafka topic
func (p *Producer) Publish(topic, message string) error {
	return p.PublishWithKey(topic, "", message)
}

// PublishWithKey publishes a message to a Kafka topic with a specific key
func (p *Producer) PublishWithKey(topic, key, message string) error {
	// For this implementation, we'll use a single broker configuration
	brokers := []string{"localhost:9092"} // This should come from config

	writer := p.getWriter(topic, brokers)

	kafkaMessage := kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		p.logger.WithError(err).Errorf("Failed to publish message to topic %s", topic)
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"topic": topic,
		"key":   key,
	}).Debug("Message published successfully")

	return nil
}

// PublishBatch publishes multiple messages to a Kafka topic
func (p *Producer) PublishBatch(topic string, messages []kafka.Message) error {
	brokers := []string{"localhost:9092"} // This should come from config

	writer := p.getWriter(topic, brokers)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.logger.WithError(err).Errorf("Failed to publish batch messages to topic %s", topic)
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"topic":        topic,
		"message_count": len(messages),
	}).Debug("Batch messages published successfully")

	return nil
}

// Close closes all Kafka writers
func (p *Producer) Close() error {
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			p.logger.WithError(err).Errorf("Failed to close writer for topic %s", topic)
		}
	}
	return nil
}

// Consumer handles Kafka message consumption
type Consumer struct {
	reader *kafka.Reader
	logger *logrus.Logger
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string, logger *logrus.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		ErrorLogger:    kafka.LoggerFunc(logger.Errorf),
	})

	return &Consumer{
		reader: reader,
		logger: logger,
	}
}

// ReadMessage reads a single message from Kafka
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// CreateTopics creates the required Kafka topics
func CreateTopics(brokers []string, topics []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		}
	}

	return conn.CreateTopics(topicConfigs...)
}