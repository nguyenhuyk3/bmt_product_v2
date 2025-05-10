package messagebroker

import (
	"bmt_product_service/internal/injectors"
	"log"
)

func InitFilmUploadConsummer() {
	filmUploadConsummer, err := injectors.InitProductUploadReader()
	if err != nil {
		log.Fatalf("an error occur when initiallizating PRODUCT UPLOAD READERs: %v\n", err)
	}

	filmUploadConsummer.InitReaders()
}
