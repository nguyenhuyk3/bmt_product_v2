package writers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var (
	brokerAddresses []string
)

func (k *KafkaWriter) getAvailableBroker() string {
	for _, broker := range brokerAddresses {
		conn, err := kafka.Dial("tcp", broker)
		if err == nil {
			conn.Close()
			return broker
		}
	}

	// log.Fatal("no available Kafka broker")
	// global.Logger.Info("no available Kafka broker")

	return ""
}

// Check if a topic already exists on the Kafka broker, and if not, automatically create it
func (k *KafkaWriter) ensureTopicExists(topic string) error {
	broker := k.getAvailableBroker()
	if broker == "" {
		return fmt.Errorf("no available Kafka broker")
	}

	// Connect to Kafka broker
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to connect to broker %s: %w", broker, err)
	}
	defer conn.Close()

	// Check for topic existence
	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		// Topic is exists
		return nil
	}

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to connect to controller broker: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}

	return nil
}
