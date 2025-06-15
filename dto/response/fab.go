package response

type FABItem struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ImageUrl string `json:"image_url"`
	Price    int32  `json:"price"`
}
