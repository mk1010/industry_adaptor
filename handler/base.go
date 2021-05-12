package handler

import (
	"os/exec"

	"github.com/apache/dubbo-go/common/logger"
)

func Init() error {
	return StartPipe("test")
}

func StartPipe(processName string) error {
	cmd := exec.Command("./"+processName, "")
	cmdStdinPipe, _ := cmd.StdinPipe()
	cmdStdoutPipe, _ := cmd.StdoutPipe()
	go func() {
		logger.Info("cmd 准备启动")
		err := cmd.Start()
		if err != nil {
			logger.Error("cmd 启动失败", err)
		}
		err = cmd.Wait()
		if err != nil {
			logger.Error("cmd 退出异常", err)
		}
		logger.Error("cmd 退出", err)
	}()
	return nil
}
