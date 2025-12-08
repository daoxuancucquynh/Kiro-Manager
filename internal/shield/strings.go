package shield

// 預編碼的敏感字串（編譯時已編碼）
// 使用 XOR key 0x5A 編碼
var (
	// 命令名稱
	CmdReg      []byte // "reg"
	CmdTaskList []byte // "tasklist"
	CmdTaskKill []byte // "taskkill"

	// reg 參數
	ArgQuery []byte // "query"
	ArgV     []byte // "/v"

	// tasklist 參數
	ArgFI        []byte // "/FI"
	ArgImageName []byte // "IMAGENAME eq Kiro.exe"
	ArgFO        []byte // "/FO"
	ArgCSV       []byte // "CSV"
	ArgNH        []byte // "/NH"

	// taskkill 參數
	ArgPID []byte // "/PID"
	ArgF   []byte // "/F"

	// Registry 路徑和值
	RegPath  []byte // `HKLM\SOFTWARE\Microsoft\Cryptography`
	RegValue []byte // "MachineGuid"
)

// init 函數在程式啟動時初始化所有編碼字串
func init() {
	codec := NewCodec()

	// 命令名稱
	CmdReg = codec.Encode("reg")
	CmdTaskList = codec.Encode("tasklist")
	CmdTaskKill = codec.Encode("taskkill")

	// reg 參數
	ArgQuery = codec.Encode("query")
	ArgV = codec.Encode("/v")

	// tasklist 參數
	ArgFI = codec.Encode("/FI")
	ArgImageName = codec.Encode("IMAGENAME eq Kiro.exe")
	ArgFO = codec.Encode("/FO")
	ArgCSV = codec.Encode("CSV")
	ArgNH = codec.Encode("/NH")

	// taskkill 參數
	ArgPID = codec.Encode("/PID")
	ArgF = codec.Encode("/F")

	// Registry 路徑和值
	RegPath = codec.Encode(`HKLM\SOFTWARE\Microsoft\Cryptography`)
	RegValue = codec.Encode("MachineGuid")
}
