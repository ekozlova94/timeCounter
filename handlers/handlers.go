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
		err := stateRepo.Save(&s)
		if err != nil {
			c.JSON(500, "Не удалось установить начало рабочего дня")
		}
		return
	}
	c.JSON(500, "Начало рабочего дня уже установлено")
}

func BreakStart(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := stateRepo.GetByDate(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	if result != nil && result.BreakStartTime == 0 {
		result.BreakStartTime = currentTime
		err := stateRepo.Save(result)
		if err != nil {
			c.JSON(500, "Не удалось установить начало перерыва")
		}
		return
	}
	if result != nil && result.BreakStartTime != 0 {
		c.JSON(500, "Начало перерыва уже установлено")
		return
	}
	c.JSON(500, "Начало рабочего дня не установлено")
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

func BreakStop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := stateRepo.GetByDate(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	if result == nil {
		c.JSON(500, "Начало рабочего дня не установлено")
		return
	}
	result.BreakStopTime = currentTime
	err = stateRepo.Save(result)
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
	sts, err = selectData(dateFromUnix, dateToUnix)
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
	result, err := dateEdit(c.Query("date"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	if result == nil {
		c.JSON(500, "Запись не найдена")
		return
	}
	startTime, err1 := strconv.Atoi(c.Query("startTime")) // преобразование из string в int
	stopTime, err2 := strconv.Atoi(c.Query("stopTime"))
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

func EditBreak(c *gin.Context) {
	result, err := dateEdit(c.Query("date"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	if result == nil {
		c.JSON(500, "Запись не найдена")
		return
	}
	breakStartTime, err1 := strconv.Atoi(c.Query("breakStartTime")) // преобразование из string в int
	breakStopTime, err2 := strconv.Atoi(c.Query("breakStopTime"))
	if err1 != nil && err2 != nil {
		c.JSON(500, "Не указано ни начало, ни окончание перерыва. Редактирование невозможно")
		return
	}
	if err1 == nil {
		result.BreakStartTime = int64(breakStartTime)
	}
	if err2 == nil {
		result.BreakStopTime = int64(breakStopTime)
	}
	err = stateRepo.Save(result)
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

func selectData(dateFromUnix int64, dateToUnix int64) ([]*models.State, error) {
	if dateFromUnix != -1 && dateToUnix != -1 {
		return stateRepo.GetByDateFromTo(dateFromUnix, dateToUnix)
	}
	if dateToUnix != -1 {
		return stateRepo.GetByDateTo(dateToUnix)
	}
	if dateFromUnix != -1 {
		return stateRepo.GetByDateFrom(dateFromUnix)
	}
	return stateRepo.GetAll()
}

func dateEdit(date string) (*models.State, error) {
	var err error
	if date == "" {
		return nil, err
	}
	dateUnix, err := converting(date)
	if err != nil {
		return nil, err
	}
	return stateRepo.GetByDate(dateUnix)
}
