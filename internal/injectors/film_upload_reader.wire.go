//go:build wireinject

package injectors

import (
	"bmt_product_service/internal/implementations/message_broker/readers"
	"bmt_product_service/internal/injectors/provider"

	"github.com/google/wire"
)

func InitFilmUploadReader() (*readers.FilmUploadReader, error) {
	wire.Build(
		uploadServiceSet,
		kafkaWriterSet,

		provider.ProvideQueries,
		readers.NewFilmUploadReader,
	)

	return &readers.FilmUploadReader{}, nil
}
