package zlog

var defaultLog *Log

func init() {
	defaultLog = New(name)
}
