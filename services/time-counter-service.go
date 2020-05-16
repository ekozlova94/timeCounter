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
			return nil, errors.New("failed to set start of work day")
		}
		return &s, nil
	}
	return nil, errors.New("the beginning of the working day is already set")
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
			return nil, errors.New("failed to set break start")
		}
		return result, nil
	}
	if result != nil && result.BreakStartTime != 0 {
		return nil, errors.New("break start already set")
	}
	return nil, errors.New("the beginning of the working day is not found")
}

func (t TimeCounterService) Stop(currentTime int64) (*models.State, error) {
	result, err := t.Repo.GetByDate(currentTime)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("the beginning of the working day is not found")
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
		return nil, errors.New("the beginning of the working day is not found")
	}
	if result.BreakStartTime == 0 {
		return nil, errors.New("break start not found")
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
		return nil, errors.New("no entries for today")
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
		return nil, errors.New("no parameters specified")
	}
	if startTime != 0 && startTime >= date && startTime < (date+86400) {
		result.StartTime = startTime
	}
	if stopTime != 0 && stopTime >= date && stopTime < (date+86400) {
		if stopTime <= startTime && startTime != 0 {
			return nil, errors.New("the end of the working day cannot be earlier than the beginning")
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
		return nil, errors.New("no parameters specified")
	}
	if breakStartTime != 0 && breakStartTime >= date && breakStartTime <= (date+86400) {
		result.BreakStartTime = breakStartTime
	}
	if breakStopTime != 0 && breakStopTime >= date && breakStopTime <= (date+86400) {
		if breakStopTime <= breakStartTime && breakStartTime != 0 {
			return nil, errors.New("the end of the working day cannot be earlier than the beginning")
		}
		result.BreakStopTime = breakStopTime
	}
	return result, t.Repo.Save(result)
}

func (t TimeCounterService) States() ([]*models.State, error) {
	return t.Repo.GetAll()
}
