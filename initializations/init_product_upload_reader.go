package initializations

import (
	"bmt_product_service/internal/injectors"
	"log"
)

func initProductUploadReader() {
	filmUploadConsummer, err := injectors.InitProductUploadReader()
	if err != nil {
		log.Fatalf("an error occur when initiallizating PRODUCT UPLOAD READERS: %v\n", err)
	}

	filmUploadConsummer.InitReaders()
}
