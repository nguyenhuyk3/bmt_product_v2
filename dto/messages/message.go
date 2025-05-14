package messages

type ReturnedObjectKeyMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ObjectKey string `json:"object_key" binding:"required"`
}

type NewFilmCreationTopic struct {
	FilmId int32 `json:"film_id" binding:"required"`
}
