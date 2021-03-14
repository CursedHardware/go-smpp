package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func runProgramWithEvent(message *Payload) {
	mutex.Lock()
	defer mutex.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, configure.Hook)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if stdin, err := cmd.StdinPipe(); err != nil {
		log.Fatalln(err)
	} else {
		go func() {
			defer stdin.Close()
			_ = json.NewEncoder(stdin).Encode(message)
		}()
	}
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

//goland:noinspection GoUnhandledErrorResult
func runProgramWithStream() func(*Payload) {
	cmd := exec.Command(configure.Hook)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		if err = cmd.Run(); err != nil {
			log.Fatalln(err)
		}
	}()
	return func(message *Payload) {
		mutex.Lock()
		defer mutex.Unlock()
		_ = json.NewEncoder(stdin).Encode(message)
		_, _ = fmt.Fprintln(stdin)
	}
}
