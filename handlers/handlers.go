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
	"timeCounter/services"
)

var Test services.TimeCounterService

func init() {
	db, err := sql.Open("sqlite3", "./db.sqlite?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	stateRepo := repositories.StateRepoImpl{
		Db: db,
	}
	Test = services.TimeCounterService{
		Repo: stateRepo,
	}
}

func Start(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := Test.Start(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func BreakStart(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := Test.BreakStart(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func Stop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := Test.Stop(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func BreakStop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := Test.BreakStop(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func Info(c *gin.Context) {
	var dateFromUnix int64 = -1
	var dateToUnix int64 = -1
	var err error
	if dateFrom := c.Query("from"); dateFrom != "" {
		dateFromUnix, err = converting(dateFrom)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
	}
	if dateTo := c.Query("to"); dateTo != "" {
		dateToUnix, err = converting(dateTo)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
	}
	var sts []*models.State
	sts, err = Test.Info(dateFromUnix, dateToUnix)
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

func Edit(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(500, "Запись не найдена")
		return
	}
	dateUnix, err := converting(date)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	startTime, err1 := strconv.Atoi(c.Query("startTime")) // преобразование из string в int
	if err1 != nil {
		c.JSON(400, err1.Error())
		return
	}
	stopTime, err2 := strconv.Atoi(c.Query("stopTime"))
	if err2 != nil {
		c.JSON(400, err2.Error())
		return
	}
	result, err := Test.Edit(dateUnix, int64(startTime), int64(stopTime))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func EditBreak(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(500, "Запись не найдена")
		return
	}
	dateUnix, err := converting(date)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	breakStartTime, err1 := strconv.Atoi(c.Query("breakStartTime")) // преобразование из string в int
	if err1 != nil {
		c.JSON(400, err1.Error())
		return
	}
	breakStopTime, err2 := strconv.Atoi(c.Query("breakStopTime"))
	if err2 != nil {
		c.JSON(400, err2.Error())
		return
	}
	result, err := Test.EditBreak(dateUnix, int64(breakStartTime), int64(breakStopTime))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func converting(date string) (int64, error) {
	t, err := time.Parse("2006-01-02", date) // преобразование из формата 2006-01-02 в типа Time
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil // преобразование из типа Time в Unix
}
