package repositories

import (
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
	"timeCounter/models"
)

type StateRepo interface {
	GetByDate(int64) (*models.State, error)
	GetByDateFromTo(int64, int64) ([]*models.State, error)
	GetByDateFrom(int64) ([]*models.State, error)
	GetByDateTo(int64) ([]*models.State, error)
	GetAll() ([]*models.State, error)
	Save(*models.State) error
}

type StateRepoImpl struct {
	Engine *xorm.Engine
}

func (o StateRepoImpl) GetByDate(t int64) (*models.State, error) {
	var s models.State
	has, err := o.Engine.Where("date($1, 'unixepoch') = date(StartTime, 'unixepoch')", t).Get(&s)
	/*rows, err := o.Engine.Query(
		`SELECT Id, StartTime, StopTime, BreakStartTime, BreakStopTime
		 FROM stats WHERE date($1, 'unixepoch') = date(StartTime, 'unixepoch')`,
		t,
	)*/
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	/*defer func() {
		if err := o.Engine.Close(); err != nil {
			log.Println(err.Error())
		}
	}()*/
	fmt.Println(&s)
	return &s, nil
}

func (o StateRepoImpl) GetByDateFromTo(dateFrom int64, dateTo int64) ([]*models.State, error) {
	var s = make([]*models.State, 0)
	has, err := o.Engine.Where("date(StartTime, 'unixepoch') BETWEEN date($1, 'unixepoch') AND date($2, 'unixepoch')", dateFrom, dateTo).Get(&s)
	/*rows, err := o.Db.Query(
		`SELECT Id, StartTime, StopTime, BreakStartTime, BreakStopTime
		 FROM stats
		 WHERE date(StartTime, 'unixepoch') BETWEEN date($1, 'unixepoch') AND date($2, 'unixepoch')`,
		dateFrom,
		dateTo,
	)*/
	if err != nil {
		return nil, err
	}
	if has {
		return nil, err
	}
	return s, nil
}

func (o StateRepoImpl) GetByDateFrom(dateFrom int64) ([]*models.State, error) {
	var s = make([]*models.State, 0)
	has, err := o.Engine.Where("date(StartTime, 'unixepoch') >= date($1, 'unixepoch')", dateFrom).Get(&s)
	/*rows, err := o.Db.Query(
			`SELECT Id, StartTime, StopTime, BreakStartTime, BreakStopTime
	 		 FROM stats WHERE date(StartTime, 'unixepoch') >= date($1, 'unixepoch')`,
			dateFrom,
		)*/
	if err != nil {
		return nil, err
	}
	if has {
		return nil, err
	}
	return s, nil
}

func (o StateRepoImpl) GetByDateTo(dateTo int64) ([]*models.State, error) {
	var s = make([]*models.State, 0)
	has, err := o.Engine.Where("date(StartTime, 'unixepoch') >= date($1, 'unixepoch')", dateTo).Get(&s)
	/*rows, err := o.Db.Query(
		`SELECT Id, StartTime, StopTime, BreakStartTime, BreakStopTime
		 FROM stats WHERE date(StartTime, 'unixepoch') >= date($1, 'unixepoch')`,
		dateTo,
	)*/
	if err != nil {
		return nil, err
	}
	if has {
		return nil, err
	}
	return s, nil
}

func (o StateRepoImpl) GetAll() ([]*models.State, error) {
	var s = make([]*models.State, 0)
	err := o.Engine.Find(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (o StateRepoImpl) Save(r *models.State) error {
	fmt.Println(r.Id)
	if r.Id > 0 {
		rowsAffected, err := o.Engine.ID(r.Id).AllCols().Update(r)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New("the number of rowsAffected is 0")
		}
		return nil
	}
	affected, err := o.Engine.Insert(r)
	if affected == 0 {
		return errors.New("the number of affected is 0")
	}
	if err != nil {
		return err
	}
	return nil
	/*	if r.Id == 0 {
		_, err := o.Db.Exec(
			`INSERT INTO stats (StartTime, StopTime, BreakStartTime, BreakStopTime) VALUES ($1, $2, $3, $4)`,
			r.StartTime,
			r.StopTime,
			r.BreakStartTime,
			r.BreakStopTime,
		)
		if err != nil {
			return err
		}
		return nil
	}*/

	/*	_, err := o.Db.Exec(
			`UPDATE stats SET StartTime=$1, StopTime=$2, BreakStartTime=$3, BreakStopTime=$4 WHERE ID=$5`,
			r.StartTime,
			r.StopTime,
			r.BreakStartTime,
			r.BreakStopTime,
			r.Id,
		)
		if err != nil {
			return err
		}
		return nil*/
}
