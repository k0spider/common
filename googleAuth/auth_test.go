package googleAuth

import (
	"fmt"
	"testing"
	"time"
)

func TestVerifyCode(t *testing.T) {
	fmt.Println(time.Now().Unix())
	//1682481138
	fmt.Println(VerifyCode("CB7X3OLJGPBFKZDRSFRQOCSZMUUSK4QC", "483168"))
}
