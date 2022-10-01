package main

import "tg/wire-di/event"

func main() {
	e := event.InitializeEvent()

	e.Start()
}
