package main

import (
	"os"
	"os/signal"
	"syscall"
)

func sigChld() {
	var sigs = make(chan os.Signal, 10) // TODO(miek): buffered channel to fix races?
	signal.Notify(sigs, syscall.SIGCHLD)

	for {
		select {
		case <-sigs:
			go reap()
		}
	}
}

func reap() {
	var wstatus syscall.WaitStatus

	pid, err := syscall.Wait4(-1, &wstatus, 0, nil)
	switch err {
	case syscall.EINTR:
		pid, err = syscall.Wait4(-1, &wstatus, 0, nil)
	case syscall.ECHILD:
		return
	}
	logPrintf("pid %d, finished, wstatus: %+v", pid, wstatus)
}