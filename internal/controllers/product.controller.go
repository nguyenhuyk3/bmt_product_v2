package controllers

import (
	"bmt_product_service/dto/request"
	"bmt_product_service/global"
	"bmt_product_service/internal/responses"
	"bmt_product_service/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	FilmService services.IFilm
	FABService  services.IFoodAndBeverage
}

func NewProductController(
	filmService services.IFilm,
	fABService services.IFoodAndBeverage) *ProductController {
	return &ProductController{
		FilmService: filmService,
		FABService:  fABService,
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

func (pc *ProductController) AddFAB(c *gin.Context) {
	name := c.PostForm("name")
	fABType := c.PostForm("type")
	image, err := c.FormFile("image")
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "no image is uploaded")
		return
	}

	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid price %s", c.PostForm("price")))
		return
	}

	req := request.AddFABReq{
		Name:  name,
		Type:  fABType,
		Image: image,
		Price: price,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, err := pc.FABService.AddFAB(ctx, req)

	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "add fab perform successfully", nil)
}

func (pc *ProductController) UpdateFAB(c *gin.Context) {
	fABId, err := strconv.Atoi(c.PostForm("fab_id"))
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid fab id %s", c.PostForm("fab_id")))
		return
	}
	name := c.PostForm("name")
	fABType := c.PostForm("type")
	image, err := c.FormFile("image")
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, "no image is uploaded")
		return
	}

	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid price %s", c.PostForm("price")))
		return
	}

	req := request.UpdateFABReq{
		FABId: int32(fABId),
		Name:  name,
		Type:  fABType,
		Image: image,
		Price: price,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, err := pc.FABService.UpdateFAB(ctx, req)

	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "update fab perform successfully", nil)
}
