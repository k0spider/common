package excel

import (
	"context"
	"testing"
)

func TestOpenFile(t *testing.T) {
	execls := NewExecl(context.Background(), "xxx.xlsx")
	c := execls.GetWriteChan()
	for i := 0; i < 10; i++ {
		data := RowData{
			Index: i,
			Data:  []interface{}{"asdasdasd", 123123},
		}
		//data = []interface{}{i, "asdcsadc", "sadcasdc"}
		c <- data
	}
	_ = execls.Close()
}
