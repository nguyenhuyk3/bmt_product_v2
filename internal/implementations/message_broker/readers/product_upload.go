package readers

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/global"
	"bmt_product_service/internal/services"
	"context"
	"log"
)

type ProductUploadReader struct {
	UploadService services.IUpload
	Writer        services.IMessageBrokerWriter
	SqlQuery      sqlc.Querier
	Context       context.Context
}

var topics = []string{
	global.RETURNED_FILM_IMAGE_OBJECT_KEY_TOPIC,
	global.RETURNED_FILM_VIDEO_OBJECT_KEY_TOPIC,
	global.RETURNED_FAB_IMAGE_OBJECT_KEY_TOPIC,
}

func NewProductUploadReader(
	uploadService services.IUpload,
	writer services.IMessageBrokerWriter,
	sqlQuery *sqlc.Queries,
) *ProductUploadReader {
	return &ProductUploadReader{
		UploadService: uploadService,
		Writer:        writer,
		SqlQuery:      sqlQuery,
		Context:       context.Background(),
	}
}

func (p *ProductUploadReader) InitReaders() {
	log.Printf("=============== Product Service is listening for film uploading messages... ===============\n\n\n")

	for _, topic := range topics {
		go p.startReader(topic)
	}
}
