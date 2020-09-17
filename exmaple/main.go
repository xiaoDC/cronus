package main

import (
	// "fmt"
	"github.com/xiaoDC/cronus/timewheel"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	tw := timewheel.New()
	tw.Start()

	for i := 1; i < 30; i++ {
		// str := fmt.Sprintf("*/%d * * * * *", i)
		// fmt.Println("----------> %s", str)
		// tw.AddCronTask(str, strconv.Itoa(i), i, true)
		tw.AddCronTask("* * * * * *", strconv.Itoa(i), i, true)
	}

	<-sigChan
}
