package controllers

import (
	"bmt_product_service/dto/request"
	"bmt_product_service/global"
	"bmt_product_service/internal/responses"
	"bmt_product_service/internal/services"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	FilmService services.IFilm
}

func NewProductController(filmService services.IFilm) *ProductController {
	return &ProductController{
		FilmService: filmService,
	}
}

func (pc *ProductController) AddFilm(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	releaseDate := c.PostForm("release_date")
	duration := c.PostForm("duration")
	genresJson := c.PostForm("genres")

	var genres []string
	if err := json.Unmarshal([]byte(genresJson), &genres); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "genres invalid format")
		return
	}

	poster, err := c.FormFile("poster")
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "no image is uploaded")
		return
	}

	trailer, err := c.FormFile("trailer")
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "no video is uploaded")
		return
	}

	req := request.AddFilmReq{
		FilmInformation: request.FilmInformation{
			Title:       title,
			Description: description,
			ReleaseDate: releaseDate,
			Genres:      genres,
			Duration:    duration,
		},
		ChangedBy: c.GetString(global.X_USER_EMAIL),
		OtherFilmInformation: request.OtherFilmInformation{
			PosterFile:  poster,
			TrailerFile: trailer,
		},
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, err := pc.FilmService.AddFilm(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "add film perform successfully", nil)
}

func (pc *ProductController) GetAllFilms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, films, err := pc.FilmService.GetAllFilms(ctx)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get all films perform successfully", films)
}

func (pc *ProductController) UpdateFilm(c *gin.Context) {
	filmId := c.PostForm("film_id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	releaseDate := c.PostForm("release_date")
	duration := c.PostForm("duration")
	genresJson := c.PostForm("genres")

	var genres []string
	if err := json.Unmarshal([]byte(genresJson), &genres); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "genres invalid format")
		return
	}

	req := request.UpdateFilmReq{
		FilmId: filmId,
		FilmInformation: request.FilmInformation{
			Title:       title,
			Description: description,
			ReleaseDate: releaseDate,
			Genres:      genres,
			Duration:    duration,
		},
		ChangedBy:            c.GetString(global.X_USER_EMAIL),
		OtherFilmInformation: request.OtherFilmInformation{},
	}

	poster, err := c.FormFile("poster")
	if err == nil {
		req.OtherFilmInformation.PosterFile = poster
	}

	trailer, err := c.FormFile("trailer")
	if err == nil {
		req.OtherFilmInformation.TrailerFile = trailer
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, err := pc.FilmService.UpdateFilm(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "update film perform successfully", nil)
}
