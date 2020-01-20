package repositories

import (
	"database/sql"
	"log"
	"timeCounter/models"
)

type StateRepo interface {
	GetByDate(int64) (*models.State, error)
	GetByDateFromTo(int64, int64) ([]*models.State, error)
	GetAll() ([]*models.State, error)
	Save(*models.State) error
}

type StateRepoImpl struct {
	Db *sql.DB
}

func (o StateRepoImpl) GetByDate(t int64) (*models.State, error) {
	rows, err := o.Db.Query("SELECT Id, StartTime, StopTime FROM stats WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')", t)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	if !rows.Next() {
		return nil, err
	}
	var s models.State
	err = rows.Scan(&s.Id, &s.StartTime, &s.StopTime)
	if err != nil {
		return nil, err
	}
	return &s, err
}

func (o StateRepoImpl) GetByDateFromTo(dateFrom int64, dateTo int64) ([]*models.State, error) {
	rows, err := o.Db.Query("SELECT Id, StartTime, StopTime FROM stats WHERE date(StartTime, 'unixepoch') BETWEEN date($1, 'unixepoch') AND date($2, 'unixepoch');", dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	return extractDataFromRows(rows)
}

func (o StateRepoImpl) GetAll() ([]*models.State, error) {
	rows, err := o.Db.Query("SELECT ID, StartTime, StopTime FROM stats")
	if err != nil {
		return nil, err
	}
	return extractDataFromRows(rows)
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

func extractDataFromRows(rows *sql.Rows) ([]*models.State, error) {
	var sts = make([]*models.State, 0)
	for rows.Next() {
		var st models.State
		err := rows.Scan(&st.Id, &st.StartTime, &st.StopTime)
		if err != nil {
			return nil, err
		}
		sts = append(sts, &st)
	}
	return sts, nil
}

// ахаха
