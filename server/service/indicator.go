package service

import (
	"log"

	"github.com/developerasun/SignalDash/server/repository"
)

type IndicatorService struct {
	repository *repository.IndicatorRepository
}

func NewIndicatorService(repo *repository.IndicatorRepository) *IndicatorService {
	indicator, err := repo.FindOne("dxy")
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("NewIndicatorService ID: %d", indicator.ID)

	return &IndicatorService{
		repository: repo,
	}
}
