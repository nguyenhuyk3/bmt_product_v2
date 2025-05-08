package readers

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/global"
	"bmt_product_service/internal/services"
	"context"
	"log"
)

type FilmUploadReader struct {
	UploadService services.IUpload
	Writer        services.IMessageBrokerWriter
	SqlQuery      sqlc.Querier
	Context       context.Context
}

var topics = []string{
	global.RETURNED_IMAGE_OBJECT_KEY_TOPIC,
	global.RETURNED_VIDEO_OBJECT_KEY_TOPIC,
}

func NewFilmUploadReader(
	uploadService services.IUpload,
	writer services.IMessageBrokerWriter,
	sqlQuery *sqlc.Queries,
) *FilmUploadReader {
	return &FilmUploadReader{
		UploadService: uploadService,
		Writer:        writer,
		SqlQuery:      sqlQuery,
		Context:       context.Background(),
	}
}

func (f *FilmUploadReader) InitReaders() {
	log.Printf("=============== Product Service is listening for film uploading messages... ===============\n\n\n")

	for _, topic := range topics {
		go f.startReader(topic)
	}
}
