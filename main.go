package main

import (
	"fmt"
	"os"
	"skynet-service/app/config"
	"skynet-service/app/exec"
	"strconv"
	"syscall"

	_ "skynet-service/spiders"
)

func main () {
	if SignalProcess() {
		fmt.Printf("skynet service is running\n")
		return
	}

	exec.Run()
}


func SignalProcess () bool {
	fp, err := os.OpenFile(config.PidFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, os.ModePerm)
	if nil != err {
		fmt.Printf("open file error: %v\n", err)
		return true
	}
	defer fp.Close()

	err = syscall.Flock(int(fp.Fd()), syscall.LOCK_EX | syscall.LOCK_NB)
	if nil != err {
		fmt.Printf("lock file error:%v\n", err)
		return true
	}

	_, err = fp.WriteString(strconv.Itoa(os.Getpid()) + "\n")
	if nil != err {
		fmt.Printf("write pid error, process already exists")
		return true
	}

	return false
}