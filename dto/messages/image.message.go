package messages

type UploadFilmImageMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ImageUrl  string `json:"image_url" binding:"required"`
}
