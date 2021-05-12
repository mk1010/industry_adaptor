package handler

import (
	"context"
	"io"
	"os/exec"

	"github.com/apache/dubbo-go/common/logger"
)

func Init(ctx context.Context) error {
	processName := "test"
	in, out, err := StartPipe(processName)
	if err != nil {
		return err
	}
	testProc := &NCLinkCommonProcess{
		ProcessName:   processName,
		CmdStdinPipe:  in,
		CmdStdoutPipe: out,
	}
	return testProc.Start(ctx)
}

func StartPipe(processName string) (io.WriteCloser, io.ReadCloser, error) {
	cmd := exec.Command("./"+processName, "")
	cmdStdinPipe, _ := cmd.StdinPipe()
	cmdStdoutPipe, _ := cmd.StdoutPipe()
	logger.Info(processName, " 准备启动")
	err := cmd.Start()
	if err != nil {
		logger.Error(processName, " 启动失败", err)
		return nil, nil, err
	}
	go func() {
		err = cmd.Wait()
		if err != nil {
			logger.Error(processName, " 退出异常", err)
		}
		logger.Error(processName, " 退出", err)
	}()
	return cmdStdinPipe, cmdStdoutPipe, nil
}
