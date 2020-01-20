package forms

import (
	"time"
	"timeCounter/models"
)

type InfoResponseForm struct {
	ID        int    `json:"id"`
	Date      string `json:"date"`
	StartTime int64  `json:"startTime"`
	StopTime  int64  `json:"stopTime"`
}

func NewInfoResponseForm(c *models.State) *InfoResponseForm {
	return &InfoResponseForm{
		ID:        c.Id,
		Date:      time.Unix(c.StartTime, 0).Format("2006-01-02"),
		StartTime: c.StartTime,
		StopTime:  c.StopTime,
	}
}
