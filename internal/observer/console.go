package observer

import "fmt"

type Console struct {
}

func (c Console) Update(msg ...any) {
	fmt.Println(msg...)
}
