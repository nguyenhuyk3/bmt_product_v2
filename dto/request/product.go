package request

import "mime/multipart"

type FilmChange struct {
	ChangedBy string
}

type OtherFilmInformation struct {
	PosterFile  *multipart.FileHeader
	TrailerFile *multipart.FileHeader
}

type FilmInformation struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	ReleaseDate string   `json:"release_date" binding:"required"`
	Genres      []string `json:"genres" binding:"required"`
	// This prop will have format as hh:mm:ss
	// When using api, we will use 2h39m
	// When stroring at databse then this will be at 02:39:00
	Duration string `json:"duration" binding:"required"`
}

type AddFilmReq struct {
	FilmInformation      FilmInformation `json:"film_information" binding:"required"`
	ChangedBy            string
	OtherFilmInformation OtherFilmInformation
}

type UpdateFilmReq struct {
	FilmId               string          `json:"film_id" binding:"required"`
	FilmInformation      FilmInformation `json:"film_information" binding:"required"`
	ChangedBy            string
	OtherFilmInformation OtherFilmInformation
}

type GetFilmByIdReq struct {
	FilmId int `json:"film_id" binding:"required"`
}

type UploadImageReq struct {
	ProductId string
	Image     *multipart.FileHeader
}

type UploadVideoReq struct {
	ProductId string
	Video     *multipart.FileHeader
}
