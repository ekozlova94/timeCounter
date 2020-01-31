package services

import (
	"errors"
	"timeCounter/models"
	"timeCounter/repositories"
)

type TimeCounterService struct {
	Repo repositories.StateRepo
}

func (t TimeCounterService) Start(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result == nil {
		var s models.State
		s.StartTime = currentTime
		err := t.Repo.Save(&s)
		if err != nil {
			return nil, errors.New("не удалось установить начало рабочего дня")
		}
		return &s, nil
	}
	return nil, errors.New("начало рабочего дня уже установлено")
}

func (t TimeCounterService) BreakStart(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result != nil && result.BreakStartTime == 0 {
		result.BreakStartTime = currentTime
		err := t.Repo.Save(result)
		if err != nil {
			return nil, errors.New("не удалось установить начало перерыва")
		}
		return result, nil
	}
	if result != nil && result.BreakStartTime != 0 {
		return nil, errors.New("начало перерыва уже установлено")
	}
	return nil, errors.New("начало рабочего дня не установлено")
}

func (t TimeCounterService) Stop(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("начало рабочего дня не установлено")
	}
	result.StopTime = currentTime
	err = t.Repo.Save(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t TimeCounterService) BreakStop(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("начало рабочего дня не установлено")
	}
	result.BreakStopTime = currentTime
	err = t.Repo.Save(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t TimeCounterService) Today(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("записи на сегодня нет")
	}
	return result, nil
}

func (t TimeCounterService) Info(dateFromUnix int64, dateToUnix int64) ([]*models.State, error) {
	if dateFromUnix != -1 && dateToUnix != -1 {
		return t.Repo.GetByDateFromTo(dateFromUnix, dateToUnix)
	}
	if dateToUnix != -1 {
		return t.Repo.GetByDateTo(dateToUnix)
	}
	if dateFromUnix != -1 {
		return t.Repo.GetByDateFrom(dateFromUnix)
	}
	return t.Repo.GetAll()
}

func (t TimeCounterService) Edit(date int64, startTime int64, stopTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(date)
	if err != nil {
		return nil, err
	}
	if startTime == 0 && stopTime == 0 {
		return nil, errors.New("не указано ни одного параметра")
	}
	if startTime != 0 && startTime >= date && startTime < (date+86400) {
		result.StartTime = startTime
	}
	if stopTime != 0 && stopTime >= date && stopTime < (date+86400) {
		if stopTime <= startTime && startTime != 0 {
			return nil, errors.New("окончание рабочего дня не может быть раньше начала")
		}
		result.StopTime = stopTime
	}
	return result, t.Repo.Save(result)
}

func (t TimeCounterService) EditBreak(date int64, breakStartTime int64, breakStopTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(date)
	if err != nil {
		return nil, err
	}
	if breakStartTime == 0 && breakStopTime == 0 {
		return nil, errors.New("не указано ни одного параметра")
	}
	if breakStartTime != 0 && breakStartTime >= date && breakStartTime <= (date+86400) {
		result.BreakStartTime = breakStartTime
	}
	if breakStopTime != 0 && breakStopTime >= date && breakStopTime <= (date+86400) {
		if breakStopTime <= breakStartTime && breakStartTime != 0 {
			return nil, errors.New("окончание рабочего дня не может быть раньше начала")
		}
		result.BreakStopTime = breakStopTime
	}
	return result, t.Repo.Save(result)
}

func (t TimeCounterService) Export() ([]*models.State, error) {
	return t.Repo.GetAll()
}
