package readers

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/dto/messages"
	"bmt_product_service/global"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"
)

func (p *ProductUploadReader) startReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_2,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_3,
		},
		GroupID:        global.PRODUCT_SERVICE_GROUP,
		Topic:          topic,
		CommitInterval: time.Second * 5,
	})
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			continue
		}

		p.processMessage(topic, message.Value)
	}
}

func (p *ProductUploadReader) processMessage(topic string, value []byte) {
	switch topic {
	case global.RETURNED_FILM_IMAGE_OBJECT_KEY_TOPIC:
		var message messages.ReturnedObjectKeyMessage
		if err := json.Unmarshal(value, &message); err != nil {
			log.Printf("failed to unmarshal image message: %v\n", err)
			return
		}

		p.handleImageObjectKeyTopic(message)

	case global.RETURNED_FILM_VIDEO_OBJECT_KEY_TOPIC:
		var message messages.ReturnedObjectKeyMessage
		if err := json.Unmarshal(value, &message); err != nil {
			log.Printf("failed to unmarshal image message: %v\n", err)
			return
		}

		p.handleVideoObjectKeyTopic(message)

	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func (p *ProductUploadReader) handleImageObjectKeyTopic(message messages.ReturnedObjectKeyMessage) {
	productId, err := strconv.Atoi(message.ProductId)
	if err != nil {
		log.Printf("product id (%s) is not in correct format: %v\n", message.ProductId, err)
		return
	}

	err = p.SqlQuery.UpdatePosterUrlAndCheckStatus(p.Context, sqlc.UpdatePosterUrlAndCheckStatusParams{
		FilmID: int32(productId),
		PosterUrl: pgtype.Text{
			String: message.ObjectKey,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("failed to update poster url for film id %d: %v\n", productId, err)
	} else {
		log.Printf("update poster url for film id %d successfully\n", productId)
	}
}

func (p *ProductUploadReader) handleVideoObjectKeyTopic(message messages.ReturnedObjectKeyMessage) {
	productId, err := strconv.Atoi(message.ProductId)
	if err != nil {
		log.Printf("product id (%s) is not in correct format: %v\n", message.ProductId, err)
		return
	}

	err = p.SqlQuery.UpdateVideoUrlAndCheckStatus(p.Context, sqlc.UpdateVideoUrlAndCheckStatusParams{
		FilmID: int32(productId),
		TrailerUrl: pgtype.Text{
			String: message.ObjectKey,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("failed to update trailer url for film id %d: %v\n", productId, err)
	} else {
		log.Printf("update trailer url for film id %d successfully\n", productId)
	}
}
