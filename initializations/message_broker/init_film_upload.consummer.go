package messagebroker

import (
	"bmt_product_service/internal/injectors"
	"log"
)

func InitFilmUploadConsummer() {
	filmUploadConsummer, err := injectors.InitFilmUploadReader()
	if err != nil {
		log.Fatalf("an error occur when initiallizating FILM UPLOAD CONSUMMER: %v\n", err)
	}

	filmUploadConsummer.InitReaders()
}
