package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
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

/*result, err := db.Exec("UPDATE stats SET StopTime=$1 WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", currentTime)
if err != nil {
	c.JSON(500, err.Error())
	return
}
rowsAffected, err := result.RowsAffected()*/

func Info(c *gin.Context) {
	var sts []*models.State
	sts = stateRepo.Query()

	if len(sts) == 0 {
		c.JSON(404, "Нет данных")
		return
	}
	m := map[string]*models.State{}
	for i := 0; i < len(sts); i++ {
		//a := strconv.Itoa(i) //конвертирование из int в string
		a := time.Unix(sts[i].StartTime, 0).Format("2006-01-02")
		m[a] = sts[i]
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
		/*rowsAffectedStartTime := stateRepo.UpdateStartTime(int64(startTime))
		if rowsAffectedStartTime != 1 {
			c.JSON(500, "Значение startTime не обновлено")
			return
		}*/
		result.StartTime = int64(startTime)
	}
	if err2 == nil {
		/*rowsAffectedStopTime := stateRepo.UpdateStopTime(int64(stopTime))
		if rowsAffectedStopTime != 1 {
			c.JSON(500, "Значение stopTime не обновлено")
			return
		}*/
		result.StopTime = int64(stopTime)
	}
	err = stateRepo.Save(result)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, result)
}

/*
fmt.Println(fmtDuration(1*time.Hour + 13*time.Minute + 23*time.Second + 10*time.Millisecond))

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02dh %02dm %02ds %03dms", h, m, s, ms)
}*/
