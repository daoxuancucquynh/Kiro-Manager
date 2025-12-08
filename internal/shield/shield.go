package shield

import "sync"

// Shield 提供統一的保護殼介面
type Shield struct {
	Codec   *StringCodec
	Builder *CmdBuilder
}

var (
	// Global 全域 Shield 實例
	Global *Shield
	once   sync.Once
)

// Init 初始化全域 Shield 實例
// 應在 main() 或 startup() 中呼叫
func Init() {
	once.Do(func() {
		Global = &Shield{
			Codec:   NewCodec(),
			Builder: NewCmdBuilder(),
		}
	})
}

// GetCodec 取得編碼器
// 如果尚未初始化，會自動初始化
func GetCodec() *StringCodec {
	if Global == nil {
		Init()
	}
	return Global.Codec
}

// GetBuilder 取得命令建構器
// 如果尚未初始化，會自動初始化
func GetBuilder() *CmdBuilder {
	if Global == nil {
		Init()
	}
	return Global.Builder
}
