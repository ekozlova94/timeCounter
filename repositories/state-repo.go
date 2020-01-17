package repositories

import (
	"database/sql"
	"log"
	"timeCounter/models"
)

type StateRepo interface {
	Request(int64) *models.State
	Add(int64) error
	Update(int64) int64
	Query() []*models.State
}

type StateRepoImpl struct {
	Db *sql.DB
}

func (o StateRepoImpl) Request(t int64) *models.State {
	rows, err := o.Db.Query("SELECT StartTime, StopTime FROM stats WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", t)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	var s models.State
	rows.Scan(&s.StartTime, &s.StopTime)
	return &s
}

func (o StateRepoImpl) Add(t int64) error {
	_, err := o.Db.Exec("INSERT INTO stats (StartTime, StopTime) VALUES ($1, 0)", t)
	if err != nil {
		return err
	}
	return nil
}

func (o StateRepoImpl) Update(t int64) int64 {
	result, _ := o.Db.Exec("UPDATE stats SET StopTime=$1 WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", t)
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected
}

func (o StateRepoImpl) Query() []*models.State {
	rows, _ := o.Db.Query("SELECT StartTime, StopTime FROM stats ORDER BY StartTime")
	sts := make([]*models.State, 0)
	for rows.Next() {
		var st models.State
		err := rows.Scan(&st.StartTime, &st.StopTime)
		if err != nil {
			log.Fatal(err)
		}
		sts = append(sts, &st)
	}
	return sts
}

// ахаха
