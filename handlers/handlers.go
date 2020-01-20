package handlers

import (
	"database/sql"
	"fmt"
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
	result := stateRepo.GetByDate(currentTime)
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
	result := stateRepo.GetByDate(currentTime)
	if result == nil {
		c.JSON(500, "Начало рабочего дня не установлено")
		return
	}
	result.StopTime = currentTime
	err := stateRepo.Save(result)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, result)
}

func Info(c *gin.Context) {
	var sts []*models.State
	sts = stateRepo.Query()

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
		c.JSON(500, "Дата не найдена")
		return
	}
	t, err := time.Parse("2006-01-02", date) // преобразование из формата 2006-01-02 в типа Time
	if err != nil {
		fmt.Println(err)
	}
	dateEdit := t.Unix()                    // преобразование из типа Time в Unix
	result := stateRepo.GetByDate(dateEdit) //переменная result хранит указатель на ячейку памяти типа modules.State, где лежит результат выполнения метода request() - список найденных записей
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
	}
	c.JSON(200, result)
}
