package dglogger

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Phone string

func (ph Phone) Secret() string {
	if ph == "" {
		return ""
	}
	return string(ph) + "***"
}

// // UnmarshalText 实现 encoding.TextUnmarshaler 接口
//
//	func (ph *Phone) UnmarshalText(b []byte) error {
//		ph.Value = string(b)
//		return nil
//	}
//
// // MarshalText 实现 encoding.TextMarshaler 接口
//
//	func (ph *Phone) MarshalText() ([]byte, error) {
//		return []byte(ph.Value), nil
//	}
type Order struct {
	Id       int64  `json:"id"`
	UserName string `json:"user_name"`
	Phone    Phone  `json:"phone"`
}

func TestLoggerJson(t *testing.T) {
	o := &Order{
		Id:       1,
		UserName: "user1",
		Phone:    Phone("1234"),
	}

	content, err := Json(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(content))

	content, err = json.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(content))
}
