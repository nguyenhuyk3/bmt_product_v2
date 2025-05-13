package product

import (
	"bmt_product_service/global"
	"context"
	"fmt"
	"strconv"
)

func (f *filmService) getFilmByIdAndCache(filmId int32) error {
	film, err := f.SqlStore.GetFilmById(context.Background(), filmId)
	if err != nil {
		return fmt.Errorf("an error occur when getting film by id (add): %v", err)
	}

	err = f.RedisClient.Save(fmt.Sprintf("%s%s", global.FILM, strconv.Itoa(int(filmId))), film, thirty_days)
	if err != nil {
		return fmt.Errorf("an error occur when saving to redis: %v", err)
	}

	return nil
}
