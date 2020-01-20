package repositories

import (
	"database/sql"
	"log"
	"timeCounter/models"
)

type StateRepo interface {
	GetByDate(int64) *models.State
	Add(int64) error
	Query() []*models.State
	Save(*models.State) error
}

type StateRepoImpl struct {
	Db *sql.DB
}

func (o StateRepoImpl) Add(int64) error {
	panic("implement me")
}

func (o StateRepoImpl) GetByDate(t int64) *models.State {
	rows, err := o.Db.Query("SELECT Id, StartTime, StopTime FROM stats WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", t)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	var s models.State
	rows.Scan(&s.Id, &s.StartTime, &s.StopTime)
	return &s
}

func (o StateRepoImpl) Query() []*models.State {
	rows, _ := o.Db.Query("SELECT ID, StartTime, StopTime FROM stats")
	var sts = make([]*models.State, 0)
	for rows.Next() {
		var st models.State
		err := rows.Scan(&st.Id, &st.StartTime, &st.StopTime)
		if err != nil {
			log.Fatal(err)
		}
		sts = append(sts, &st)
	}
	return sts
}

func (o StateRepoImpl) Save(r *models.State) error {
	if r.Id == 0 {
		_, err := o.Db.Exec("INSERT INTO stats (StartTime, StopTime) VALUES ($1, $2)", r.StartTime, r.StopTime)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := o.Db.Exec("UPDATE stats SET StartTime=$1, StopTime=$2 WHERE ID=$3", r.StartTime, r.StopTime, r.Id)
	if err != nil {
		return err
	}
	return nil
}

// ахаха
