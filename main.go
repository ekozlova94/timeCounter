package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"timeCounter/models"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.POST("/api/start", start)
	r.POST("/api/stop", stop)
	r.GET("/api/info", info)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	//fmt.Println(fmtDuration(1*time.Hour + 13*time.Minute + 23*time.Second + 10*time.Millisecond))
}

func start(c *gin.Context) {
	currentTime := time.Now().Unix()
	rows, err := db.Query("SELECT * FROM stats WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", currentTime)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() != true {
		db.Exec("INSERT INTO stats (StartTime, StopTime) VALUES ($1, 0)", currentTime)
		time.Now().Date()
		c.Status(200)
		return
	}
	c.JSON(500, "Начало рабочего дня уже было установлено")

}

func stop(c *gin.Context) {
	currentTime := time.Now().Unix() //.Add(1*time.Hour)
	result, err := db.Exec("UPDATE stats SET StopTime=$1 WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", currentTime)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 {
		c.JSON(500, "Что-то пошло не так")
		return
	}
	c.Status(200)
}

func info(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM stats ORDER BY StartTime")
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer rows.Close()
	var sts []models.State
	for rows.Next() {
		var s models.State
		rows.Scan(&s.StartTime, &s.StopTime)
		sts = append(sts, s)
	}
	if len(sts) == 0 {
		c.JSON(404, "Нет данных")
		return
	}

	m := map[string]models.State{}
	for i := 0; i < len(sts); i++ {
		//a := strconv.Itoa(i) //конвертирование из int в string
		a := time.Unix(sts[i].StartTime, 0).Format("2006-01-02")
		m[a] = sts[i]
	}
	c.JSON(200, m)
}

/*
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
