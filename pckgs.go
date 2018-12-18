package telegram

type EventSign struct { //get data from user
	State int // 0 - email, 1 - phone, 2 - event
	Name  string
	Email string
	Phone string
	Event string
}
