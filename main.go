//go:generate swag init
package main

import (
	"fmt"
	"message-scheduler/log"
)

func main() {
	fmt.Println("Hello World")

	log.Logger.Info().Msg("Message Scheduler exiting...")
}
