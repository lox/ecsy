package cmds

import (
	"io"
	"log"
	"os"
)

type Ui struct {
	*log.Logger
	Error, Debug *log.Logger
	Exit         func(code int)
}

func (ui *Ui) DebugWriter() io.Writer {
	return &logWriter{ui.Debug}
}

func NewUi() *Ui {
	return &Ui{
		Logger: log.New(os.Stdout, "", 0),
		Error:  log.New(os.Stderr, "", 0),
		Debug:  log.New(os.Stderr, "DEBUG ", log.LstdFlags),
		Exit:   os.Exit,
	}
}

type logWriter struct{ *log.Logger }

func (w logWriter) Write(b []byte) (int, error) {
	w.Printf("%s", b)
	return len(b), nil
}
