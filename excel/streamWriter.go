package excel

import (
	"context"
	"github.com/k0spider/common/log"
	"github.com/xuri/excelize/v2"
	"sync"
)

type RowData struct {
	Index int
	Data  []interface{}
}

type execl struct {
	f            *excelize.File
	writeChannel chan RowData
	currentSheet string
	ctx          context.Context
	fileName     string
	sync.WaitGroup
	done chan bool
}

func NewExecl(ctx context.Context, fileName string) *execl {
	e := new(execl)
	e.ctx = ctx
	e.f = excelize.NewFile()
	e.fileName = fileName
	e.writeChannel = make(chan RowData)
	e.done = make(chan bool)
	e.currentSheet = "Sheet1"
	go e.receive()
	return e
}

func (e *execl) GetWriteChan() chan<- RowData {
	return e.writeChannel
}

func (e *execl) Close() error {
	close(e.writeChannel)
	<-e.done
	e.Wait()
	if err := e.f.SaveAs(e.fileName); err != nil {
		return err
	}
	return e.f.Close()
}

func (e *execl) receive() {
	i := 1
	for {
		data, ok := <-e.writeChannel
		if !ok {
			e.done <- true
			break
		}
		if data.Index == 0 {
			i++
			data.Index = i
		}
		go func(data RowData) {
			e.write(data)
		}(data)
	}
}

func (e *execl) write(data interface{}) {
	e.Add(1)
	defer e.Done()
	switch message := data.(type) {
	case RowData:
		cell, err := excelize.CoordinatesToCellName(1, message.Index)
		if err != nil {
			log.WithContext(e.ctx).Errorf("CoordinatesToCellName error:%s", err.Error())
			return
		}
		if err = e.f.SetSheetRow(e.currentSheet, cell, &message.Data); err != nil {
			log.WithContext(e.ctx).Errorf("CoordinatesToCellName message:%v error:%s", message.Data, err.Error())
		}
		return
	}
}
