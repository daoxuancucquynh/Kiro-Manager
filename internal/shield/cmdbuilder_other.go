//go:build !windows

package shield

import (
	"os/exec"
)

// SetHidden 非 Windows 平台的空實作
// 其他平台不需要隱藏視窗
func (b *CmdBuilder) SetHidden(cmd *exec.Cmd) {
	// 非 Windows 平台不需要特殊處理
}
