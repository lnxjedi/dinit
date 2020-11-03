package main

// Adapted from https://github.com/ramr/go-reaper/blob/master/reaper.go
// No license published there...

import (
	"os"
	"os/signal"
	"syscall"
)

func childSignal(notify chan bool) {
	var sigs = make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGCHLD)

	for {
		<-sigs
		select {
		case notify <- true:
		default:
			// Channel full, does not matter as we wait for all children.
		}
	}
}

func reap() {
	var wstatus syscall.WaitStatus
	notify := make(chan bool, 1)

	go childSignal(notify)

	for {
		pid, err := syscall.Wait4(-1, &wstatus, 0, nil)
		for err == syscall.EINTR {
			pid, err = syscall.Wait4(-1, &wstatus, 0, nil)
		}
		if err == syscall.ECHILD {
			// it's odd that we would get this and this used to 'break' the loop. Now
			// log this has happened, but keep waiting.
			lg.Printf("wait4() returned ECHILD")
			continue
		}
		lg.Printf("pid %d, finished, wstatus: %+v", pid, wstatus)
	}
}
