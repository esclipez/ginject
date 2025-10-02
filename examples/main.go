package main

import (
	"time"

	"github.com/esclipez/ginject/boot"
)

func main() {
	time.AfterFunc(1*time.Second, func() {
		boot.Shutdown()
	})
	boot.RunApplication()
}
