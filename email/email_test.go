package email

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	err := NewEmail(&Config{
		UserName: "111@gmail.com",
		Password: "xxxxxxx",
	}).Gmail().Send(&EmailParameter{
		To:       []string{"111@gmail.com"},
		Subject:  "test",
		TextBody: []byte("code:121212"),
	})
	fmt.Println(err)
}
