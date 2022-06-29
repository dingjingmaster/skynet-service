package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"spider/logs"
	"syscall"
	"time"
)

var runIndex = 0

type Daemon struct {
	Stop 			bool			// 是否退出
}

/**
 * 把本身程序转化为后台运行(启动一个子进程，然后自己退出)
 */
func Background (isExit bool) (*exec.Cmd, error) {
	// 判断是子进程还是父进程
	runIndex++

	// 设置子进程环境变量，启动次数
	env := os.Environ()
	env = append (env, fmt.Sprintf("%d", runIndex))

	// 启动子进程
	cmd, err := startProc (os.Args, env)
	if nil != err {
		logs.Log.Error("sub process start error:%v", err)
		return nil, err
	} else {
		logs.Log.Informational("sub process start ok!")
	}

	if isExit {
		os.Exit(0)
	}

	return cmd, nil
}

func NewDaemon () *Daemon {
	return &Daemon {}
}

/**
 * 启动后台守护进程
 */
func (d* Daemon) Run () {
	Background(true)		// 启动一个守护进程之后退出

	errNum := 0
	// 守护进程启动一个子进程并循环监视
	for !d.Stop {
		// daemon 描述信息
		logs.Log.Emergency("pid:%d, errNum:%d sub process fail", os.Getpid(), errNum)

		t := time.Now().Unix()	// 启动时间戳
		cmd, err := Background(false)
		if nil != err {
			logs.Log.Alert("sub process start failed! %s", err)
			errNum++
			continue
		}

		// 子进程
		if nil == cmd {
			logs.Log.Informational("sub process pid:%d starting...", os.Getpid())
			break
		}

		// 父进程
		err = cmd.Wait()
		dat := time.Now().Unix() - t

		errNum++

		logs.Log.Informational("sub process pid:%d runs %d second, error:%v", cmd.ProcessState.Pid(), dat, err)
	}
}

func startProc (args, env []string) (*exec.Cmd, error) {
	cmd := &exec.Cmd {
		Path:			args[0],
		Args:			args,
		Env: 			env,
		SysProcAttr: 	&syscall.SysProcAttr {Setsid: true},
	}

	err := cmd.Start()
	if nil != err {
		return nil, err
	}

	return cmd, nil
}