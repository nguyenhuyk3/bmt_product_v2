package messages

type UploadFilmVideoMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	VideoUrl  string `json:"video_url" binding:"required"`
}
