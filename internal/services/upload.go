package services

import "bmt_product_service/dto/request"

type IUpload interface {
	UploadFilmImageToS3(message request.UploadImageReq) error
	UploadFilmVideoToS3(message request.UploadVideoReq) error
}
