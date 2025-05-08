package services

import "bmt_product_service/dto/request"

type IUpload interface {
	UploadProductImageToS3(message request.UploadImageReq, productType string) error
	UploadFilmVideoToS3(message request.UploadVideoReq) error
	DeleteObject(objectURL string) error
}
