//go:build windows

package machineid

import (
	"errors"
	"strings"

	"kiro-manager/internal/shield"
)

// getWindowsMachineId 使用 reg query 命令讀取 Registry 中的 MachineGuid
// 使用 Shield 模組保護敏感字串，避免防毒軟體靜態分析誤報
func getWindowsMachineId() (string, error) {
	// 使用 Shield 建構 reg query 命令
	// 原始命令: reg query "HKLM\SOFTWARE\Microsoft\Cryptography" /v MachineGuid
	builder := shield.GetBuilder()
	cmd := builder.Build(shield.CmdReg, shield.ArgQuery, shield.RegPath, shield.ArgV, shield.RegValue)
	if cmd == nil {
		return "", errors.New("failed to build reg query command")
	}
	builder.SetHidden(cmd)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 輸出格式:
	// HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography
	//     MachineGuid    REG_SZ    xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "MachineGuid") && strings.Contains(line, "REG_SZ") {
			// 分割並取得最後一個欄位（UUID）
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return strings.ToLower(parts[len(parts)-1]), nil
			}
		}
	}

	return "", errors.New("MachineGuid not found in registry")
}

// getDarwinMachineId Windows 平台不支援
func getDarwinMachineId() (string, error) {
	return "", errors.New("Darwin-only function called on Windows")
}

// getLinuxMachineId Windows 平台不支援
func getLinuxMachineId() (string, error) {
	return "", errors.New("Linux-only function called on Windows")
}
