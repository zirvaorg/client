package observer

type Observer interface {
	Update(msg ...any)
}
