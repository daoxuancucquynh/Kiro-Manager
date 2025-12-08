package shield

import (
	"os/exec"
)

// CmdBuilder 動態建構系統命令
type CmdBuilder struct {
	codec *StringCodec
}

// NewCmdBuilder 建立新的命令建構器
func NewCmdBuilder() *CmdBuilder {
	return &CmdBuilder{
		codec: NewCodec(),
	}
}

// Build 使用編碼後的名稱和參數建構命令
// encodedName: 編碼後的命令名稱
// encodedArgs: 編碼後的參數列表
func (b *CmdBuilder) Build(encodedName []byte, encodedArgs ...[]byte) *exec.Cmd {
	if len(encodedName) == 0 {
		return nil
	}

	name := b.codec.Decode(encodedName)
	args := make([]string, len(encodedArgs))
	for i, arg := range encodedArgs {
		args[i] = b.codec.Decode(arg)
	}

	return exec.Command(name, args...)
}

// BuildWithRawArgs 使用編碼後的名稱和原始參數建構命令
// 用於需要動態參數的情況（如 PID）
func (b *CmdBuilder) BuildWithRawArgs(encodedName []byte, rawArgs ...string) *exec.Cmd {
	if len(encodedName) == 0 {
		return nil
	}

	name := b.codec.Decode(encodedName)
	return exec.Command(name, rawArgs...)
}
