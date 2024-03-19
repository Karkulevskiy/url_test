package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"url/models"
)

type Service interface {
	GetURL(context.Context) (*models.Site, error)
	CutURL(context.Context) (*models.Site, error)
}

type SiteService struct{}

func NewSiteService() *SiteService {
	return &SiteService{}
}

/* func (ss *SiteService) CutURL(context context.Context) (*models.Site, error) {

} */

func (ss *SiteService) GetURL(context context.Context) (*models.Site, error) {
	response, err := http.Get("http://localhost:8080/")

	//https://www.youtube.com/watch?v=sqj4UzN4OpU - Сделать DI при проверке на ошибки
	if err != nil {
		log.Println("Error in getting url")
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Request.Body)
	if err != nil {
		log.Println("Error while decoding body")
		return nil, err
	}
	site := &models.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		log.Println("Error while creating site")
		return nil, err
	}
	return site, nil
}
