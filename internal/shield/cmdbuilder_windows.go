//go:build windows

package shield

import (
	"os/exec"
	"syscall"
)

// SetHidden 設定命令為隱藏視窗模式（Windows）
// 使用 STARTF_USESHOWWINDOW + SW_HIDE 來隱藏視窗
func (b *CmdBuilder) SetHidden(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
