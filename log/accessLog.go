package log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

var fileWriterHub sync.Map
var _, timeZone = time.Now().Zone()

type accessLog struct {
	Path        string
	Prefix      string
	File        *os.File
	CurrentDate int64
}

func NewAccessLog(logPath string, prefix string) *accessLog {
	w := &accessLog{
		Path:   logPath,
		Prefix: prefix,
	}
	writer, ok := fileWriterHub.Load(w.getKey())
	if ok {
		return writer.(*accessLog)
	}

	fileWriterHub.Store(w.getKey(), w)
	return w
}

func (ac *accessLog) Close() {
	ac.File.Close()
	fileWriterHub.Delete(ac.getKey())
}

func (ac *accessLog) Write(p []byte) (int, error) {
	ac.checkDate()
	return ac.File.Write(p)
}

func (ac *accessLog) getKey() string {
	return path.Join(ac.Path, ac.Prefix) + "/"
}

func (ac *accessLog) checkDate() {
	var err error
	now := time.Now()
	date := (now.Unix() + int64(timeZone)) / 86400

	if date == ac.CurrentDate {
		return
	}
	ac.File.Close()
	os.MkdirAll(fmt.Sprintf("%v/%v", ac.Path, ac.Prefix), os.ModePerm)
	fileName := fmt.Sprintf("%v/%v/%v%v", ac.Path, ac.Prefix, now.Format("2006-01-02"), ".log")
	ac.File, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(fmt.Sprintf("open log file failed. err:%s, file:%s\n", err, fileName))
	}
	ac.CurrentDate = date
}
