package main

import (
	"io"
	"log"
	"os"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var DefaultUi *Ui

func init() {
	DefaultUi = NewUi()
	log.SetFlags(0)
	log.SetOutput(DefaultUi.DebugWriter())
}

type Ui struct {
	*log.Logger
	Error, Debug *log.Logger
	Exit         func(code int)
}

func (ui *Ui) DebugWriter() io.Writer {
	return &logWriter{ui.Debug}
}

func (ui *Ui) EcsEventPrinter() func(e *ecs.ServiceEvent) {
	return func(e *ecs.ServiceEvent) {
		ui.Println(*e.Message)
	}
}

func NewUi() *Ui {
	return &Ui{
		Logger: log.New(os.Stdout, "INFO ", log.LstdFlags),
		Error:  log.New(os.Stderr, "ERROR ", log.LstdFlags),
		Debug:  log.New(os.Stderr, "DEBUG ", log.LstdFlags),
		Exit:   os.Exit,
	}
}

type logWriter struct{ *log.Logger }

func (w logWriter) Write(b []byte) (int, error) {
	w.Printf("%s", b)
	return len(b), nil
}
