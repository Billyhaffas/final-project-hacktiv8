package transaction

import "time"

type Transaction struct {
	Id        string
	UserId    string
	RoomId    string
	Qty       int
	OrderDate time.Time
	Subtotal  float32
	Total     float32
}
