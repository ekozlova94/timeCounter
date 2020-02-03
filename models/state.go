package models

type State struct {
	Id             int   `xorm:"pk autoincr 'Id'"`
	StartTime      int64 `xorm:"'StartTime'"`
	StopTime       int64 `xorm:"'StopTime'"`
	BreakStartTime int64 `xorm:"'BreakStartTime'"`
	BreakStopTime  int64 `xorm:"'BreakStopTime'"`
}
