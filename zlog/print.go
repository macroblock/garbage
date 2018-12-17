package zlog

// Panic -
func (o *Log) Panic(cause interface{}, msg ...interface{}) {
	o.Log(LevelPanic, cause, msg...)
}

// Panicf -
func (o *Log) Panicf(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelPanic, cause, format, msg...)
}

// Critical -
func (o *Log) Critical(cause interface{}, msg ...interface{}) {
	o.Log(LevelCritical, cause, msg...)
}

// Criticalf -
func (o *Log) Criticalf(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelCritical, cause, format, msg...)
}

// Error -
func (o *Log) Error(cause interface{}, msg ...interface{}) {
	o.Log(LevelError, cause, msg...)
}

// Errorf -
func (o *Log) Errorf(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelError, cause, format, msg...)
}

// Warning -
func (o *Log) Warning(cause interface{}, msg ...interface{}) {
	o.Log(LevelWarning, cause, msg...)
}

// Warningf -
func (o *Log) Warningf(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelWarning, cause, format, msg...)
}

// Notice -
func (o *Log) Notice(cause interface{}, msg ...interface{}) {
	o.Log(LevelNotice, cause, msg...)
}

// Noticef -
func (o *Log) Noticef(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelNotice, cause, format, msg...)
}

// Info -
func (o *Log) Info(cause interface{}, msg ...interface{}) {
	o.Log(LevelInfo, cause, msg...)
}

// Infof -
func (o *Log) Infof(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelInfo, cause, format, msg...)
}

// Degug -
func (o *Log) Degug(cause interface{}, msg ...interface{}) {
	o.Log(LevelDebug, cause, msg...)
}

// Debugf -
func (o *Log) Debugf(cause interface{}, format string, msg ...interface{}) {
	o.Logf(LevelDebug, cause, format, msg...)
}
