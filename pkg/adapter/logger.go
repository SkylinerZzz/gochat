package adapter

type Logger struct {
	success int    // task completion count
	failure int    // task incomplete count
	name    string // logger name
}

func NewLogger(name string) Logger {
	return Logger{name: name}
}

func (l Logger) Run(err error) {
	if err != nil {
		l.failure++
	} else {
		l.success++
	}
}

func (l Logger) Name() string {
	return l.name
}
