package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
	"timeCounter/forms"
	"timeCounter/models"
	"timeCounter/repositories"
)

var stateRepo repositories.StateRepo

func init() {
	db, err := sql.Open("sqlite3", "./db.sqlite?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	stateRepo = repositories.StateRepoImpl{
		Db: db,
	}
}

func Start(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := stateRepo.GetByDate(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
	}
	if result == nil {
		var s models.State
		s.StartTime = currentTime
		s.StopTime = 0
		err := stateRepo.Save(&s)
		if err != nil {
			c.JSON(500, "Не удалось установить начало рабочего дня")
		}
		return
	}
	c.JSON(500, "Начало рабочего дня уже установлено")
}

func Stop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := stateRepo.GetByDate(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
	}
	if result == nil {
		c.JSON(500, "Начало рабочего дня не установлено")
		return
	}
	result.StopTime = currentTime
	err = stateRepo.Save(result)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, result)
}

func Info(c *gin.Context) {
	dateFrom := c.Query("from")
	dateTo := c.Query("to")

	if dateFrom != "" && dateTo != "" {
		dateFromUnix, err := converting(dateFrom)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		dateToUnix, err := converting(dateTo)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		sts, err := stateRepo.GetByDateFromTo(dateFromUnix, dateToUnix)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		if len(sts) == 0 {
			c.JSON(404, "Нет данных")
			return
		}
		var m []*forms.InfoResponseForm

		for i := 0; i < len(sts); i++ {
			m = append(m, forms.NewInfoResponseForm(sts[i]))
		}
		c.JSON(200, m)
		return
	}
	sts, err := stateRepo.GetAll()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	if len(sts) == 0 {
		c.JSON(404, "Нет данных")
		return
	}
	var m []*forms.InfoResponseForm

	for i := 0; i < len(sts); i++ {
		m = append(m, forms.NewInfoResponseForm(sts[i]))
	}
	c.JSON(200, m)
}
func converting(date string) (int64, error) {
	t, err := time.Parse("2006-01-02", date) // преобразование из формата 2006-01-02 в типа Time
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil // преобразование из типа Time в Unix
}
func Edit(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(500, "Дата не найдена")
		return
	}
	dateUnix, err := converting(date)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	result, err := stateRepo.GetByDate(dateUnix)
	if err != nil {
		c.JSON(500, err.Error())
	}
	if result == nil {
		c.JSON(500, "Запись не найдена")
		return
	}
	startTime, err1 := strconv.Atoi(c.Query("startTime")) // преобразование из string в int
	stopTime, err2 := strconv.Atoi(c.Query("stopTime"))   // преобразование из string в  int
	if err1 != nil && err2 != nil {
		c.JSON(500, "Не указано ни начало, ни окончание рабочего дня. Редактирование невозможно")
		return
	}
	if err1 == nil {
		result.StartTime = int64(startTime)
	}
	if err2 == nil {
		result.StopTime = int64(stopTime)
	}
	err = stateRepo.Save(result)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}
