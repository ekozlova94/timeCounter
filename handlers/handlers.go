package handlers

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"timeCounter/forms"
	"timeCounter/models"
	"timeCounter/repositories"
	"timeCounter/services"
	"xorm.io/core"
)

var TimeCounterService services.TimeCounterService

func init() {
	engine, err := xorm.NewEngine("sqlite3", "./db.sqlite?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMapper(core.GonicMapper{})
	/*db, err := sql.Open("sqlite3", "./db.sqlite?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}*/
	if err = engine.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = engine.Sync2(new(models.State)); err != nil {
		log.Fatal(err)
	}
	stateRepo := repositories.StateRepoImpl{
		Engine: engine,
	}
	TimeCounterService = services.TimeCounterService{
		Repo: stateRepo,
	}
}

func Start(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := TimeCounterService.Start(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func BreakStart(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := TimeCounterService.BreakStart(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func Stop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := TimeCounterService.Stop(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func BreakStop(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := TimeCounterService.BreakStop(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func Today(c *gin.Context) {
	currentTime := time.Now().Unix()
	result, err := TimeCounterService.Today(currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, forms.NewInfoResponseForm(result))
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
	sts, err = TimeCounterService.Info(dateFromUnix, dateToUnix)
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
		c.JSON(500, "record not found")
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
	result, err := TimeCounterService.Edit(dateUnix, int64(startTime), int64(stopTime))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

func EditBreak(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(500, "record not found")
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
	result, err := TimeCounterService.EditBreak(dateUnix, int64(breakStartTime), int64(breakStopTime))
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

func Export(c *gin.Context) {
	exportData, errExport := TimeCounterService.States()
	if errExport != nil {
		c.JSON(400, errExport.Error())
		return
	}
	file, errCreateTempFile := ioutil.TempFile("temp", "TempFile-*.csv")
	if errCreateTempFile != nil {
		c.JSON(400, "Unable to create temp file:")
		return
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil {
			log.Println(err.Error())
		}
	}()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err.Error())
		}
	}()
	data := "Дата" + "," + "Начало рабочего дня" + "," + "Начало перерыва" + "," + "Окончание перерыва" + "," + "Окончание рабочего дня" + "\r"
	_, errWriteString := file.WriteString(data)
	if errWriteString != nil {
		c.JSON(500, "Unable to write to file")
		return
	}
	for i := 0; i < len(exportData); i++ {
		date := time.Unix(exportData[i].StartTime, 0).Format("2006-01-02")
		startTime := time.Unix(exportData[i].StartTime, 0).Format("15:04:05")
		stopTime := time.Unix(exportData[i].StopTime, 0).Format("15:04:05")
		breakStartTime := time.Unix(exportData[i].BreakStartTime, 0).Format("15:04:05")
		breakStopTime := time.Unix(exportData[i].BreakStopTime, 0).Format("15:04:05")
		data = date + "," + startTime + "," + breakStartTime + "," + breakStopTime + "," + stopTime + "\r"
		_, errWriteString := file.WriteString(data)
		if errWriteString != nil {
			c.JSON(500, "Unable to write to file")
			return
		}
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+"Report.csv")
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.Name())
}
