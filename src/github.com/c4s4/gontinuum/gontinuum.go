package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		start := time.Now()
		config := LoadConfig(os.Args[i])
		builds := BuildModules(config)
		duration := time.Since(start)
		fmt.Println("Done in", duration)
		fmt.Println(builds.String())
		SendEmail(builds, start, duration, config.Email)
	}
}
