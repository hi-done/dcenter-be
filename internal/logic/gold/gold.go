package gold

import (
	"dcenter/internal/service"
	"dcenter/utility/gold"
	"time"
)

var goldService = gold.NewGoldService()

type sGold struct {
}

func init() {
	service.RegisterGold(New())
}

func New() service.IGold {
	return &sGold{}
}

func (s *sGold) GetGoldPrice() (float64, string, time.Time) {
	return goldService.GetGoldPrice()
}
