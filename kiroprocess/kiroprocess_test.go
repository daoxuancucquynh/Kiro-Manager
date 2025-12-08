//go:build windows

package kiroprocess

import (
	"testing"

	"kiro-manager/internal/shield"
)

// TestIntegration_GetWindowsKiroProcesses 測試 getWindowsKiroProcesses 使用 Shield 正常運作
// 驗證 Shield 的編碼字串和 CmdBuilder 能正確建構 tasklist 命令
// **Validates: Requirements 4.1**
func TestIntegration_GetWindowsKiroProcesses(t *testing.T) {
	// 確保 Shield 已初始化
	shield.Init()

	// 測試 getWindowsKiroProcesses 能正常執行
	// 即使沒有 Kiro 進程運行，也應該返回空列表而非錯誤
	processes, err := getWindowsKiroProcesses()
	if err != nil {
		t.Fatalf("getWindowsKiroProcesses() returned error: %v", err)
	}

	// 驗證返回的是有效的 slice（可能為空）
	if processes == nil {
		t.Fatal("getWindowsKiroProcesses() returned nil, expected empty slice")
	}

	// 如果有進程，驗證每個進程的結構
	for i, p := range processes {
		if p.PID <= 0 {
			t.Errorf("Process %d has invalid PID: %d", i, p.PID)
		}
		if p.Name == "" {
			t.Errorf("Process %d has empty name", i)
		}
	}

	t.Logf("Found %d Kiro processes", len(processes))
}

// TestIntegration_ShieldEncodedStringsForTasklist 測試 Shield 編碼字串用於 tasklist 命令
// 驗證 CmdTaskList 和相關參數能正確解碼
// **Validates: Requirements 4.1**
func TestIntegration_ShieldEncodedStringsForTasklist(t *testing.T) {
	shield.Init()
	codec := shield.GetCodec()

	// 驗證 tasklist 相關的編碼字串能正確解碼
	testCases := []struct {
		name     string
		encoded  []byte
		expected string
	}{
		{"CmdTaskList", shield.CmdTaskList, "tasklist"},
		{"ArgFI", shield.ArgFI, "/FI"},
		{"ArgImageName", shield.ArgImageName, "IMAGENAME eq Kiro.exe"},
		{"ArgFO", shield.ArgFO, "/FO"},
		{"ArgCSV", shield.ArgCSV, "CSV"},
		{"ArgNH", shield.ArgNH, "/NH"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded := codec.Decode(tc.encoded)
			if decoded != tc.expected {
				t.Errorf("Decode(%s) = %q, want %q", tc.name, decoded, tc.expected)
			}
		})
	}
}

// TestIntegration_ShieldEncodedStringsForTaskkill 測試 Shield 編碼字串用於 taskkill 命令
// 驗證 CmdTaskKill 和相關參數能正確解碼
// **Validates: Requirements 4.2**
func TestIntegration_ShieldEncodedStringsForTaskkill(t *testing.T) {
	shield.Init()
	codec := shield.GetCodec()

	// 驗證 taskkill 相關的編碼字串能正確解碼
	testCases := []struct {
		name     string
		encoded  []byte
		expected string
	}{
		{"CmdTaskKill", shield.CmdTaskKill, "taskkill"},
		{"ArgPID", shield.ArgPID, "/PID"},
		{"ArgF", shield.ArgF, "/F"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded := codec.Decode(tc.encoded)
			if decoded != tc.expected {
				t.Errorf("Decode(%s) = %q, want %q", tc.name, decoded, tc.expected)
			}
		})
	}
}

// TestIntegration_CmdBuilderForTasklist 測試 CmdBuilder 建構 tasklist 命令
// 驗證 Build() 能正確建構 tasklist 命令
// **Validates: Requirements 4.1**
func TestIntegration_CmdBuilderForTasklist(t *testing.T) {
	shield.Init()
	builder := shield.GetBuilder()

	// 建構 tasklist 命令
	cmd := builder.Build(shield.CmdTaskList, shield.ArgFI, shield.ArgImageName, shield.ArgFO, shield.ArgCSV, shield.ArgNH)
	if cmd == nil {
		t.Fatal("Build() returned nil")
	}

	// 驗證命令路徑包含 tasklist
	if cmd.Path == "" {
		t.Error("cmd.Path is empty")
	}

	// 驗證參數數量正確（命令名 + 5 個參數）
	expectedArgCount := 6 // tasklist, /FI, "IMAGENAME eq Kiro.exe", /FO, CSV, /NH
	if len(cmd.Args) != expectedArgCount {
		t.Errorf("cmd.Args has %d elements, want %d", len(cmd.Args), expectedArgCount)
	}

	t.Logf("Built command: %v", cmd.Args)
}

// TestIntegration_CmdBuilderForTaskkill 測試 CmdBuilder 建構 taskkill 命令
// 驗證 BuildWithRawArgs() 能正確建構 taskkill 命令
// **Validates: Requirements 4.2**
func TestIntegration_CmdBuilderForTaskkill(t *testing.T) {
	shield.Init()
	builder := shield.GetBuilder()
	codec := shield.GetCodec()

	// 建構 taskkill 命令（使用 BuildWithRawArgs 因為 PID 是動態值）
	pidArg := codec.Decode(shield.ArgPID)
	forceArg := codec.Decode(shield.ArgF)
	testPID := "12345"

	cmd := builder.BuildWithRawArgs(shield.CmdTaskKill, pidArg, testPID, forceArg)
	if cmd == nil {
		t.Fatal("BuildWithRawArgs() returned nil")
	}

	// 驗證命令路徑包含 taskkill
	if cmd.Path == "" {
		t.Error("cmd.Path is empty")
	}

	// 驗證參數數量正確（命令名 + 3 個參數）
	expectedArgCount := 4 // taskkill, /PID, 12345, /F
	if len(cmd.Args) != expectedArgCount {
		t.Errorf("cmd.Args has %d elements, want %d", len(cmd.Args), expectedArgCount)
	}

	t.Logf("Built command: %v", cmd.Args)
}

// TestIntegration_ParseTasklistOutput 測試 tasklist CSV 輸出解析
// 驗證 parseTasklistOutput 能正確解析各種輸出格式
func TestIntegration_ParseTasklistOutput(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		expectedCount  int
		expectedPID    int
		expectedName   string
	}{
		{
			name:           "Single process",
			output:         `"Kiro.exe","12345","Console","1","123,456 K"`,
			expectedCount:  1,
			expectedPID:    12345,
			expectedName:   "Kiro.exe",
		},
		{
			name:           "Multiple processes",
			output:         "\"Kiro.exe\",\"12345\",\"Console\",\"1\",\"123,456 K\"\n\"Kiro.exe\",\"67890\",\"Console\",\"1\",\"234,567 K\"",
			expectedCount:  2,
			expectedPID:    12345,
			expectedName:   "Kiro.exe",
		},
		{
			name:           "No processes (INFO message)",
			output:         "INFO: No tasks are running which match the specified criteria.",
			expectedCount:  0,
			expectedPID:    0,
			expectedName:   "",
		},
		{
			name:           "Empty output",
			output:         "",
			expectedCount:  0,
			expectedPID:    0,
			expectedName:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			processes, err := parseTasklistOutput(tc.output)
			if err != nil {
				t.Fatalf("parseTasklistOutput() returned error: %v", err)
			}

			if len(processes) != tc.expectedCount {
				t.Errorf("parseTasklistOutput() returned %d processes, want %d", len(processes), tc.expectedCount)
			}

			if tc.expectedCount > 0 && len(processes) > 0 {
				if processes[0].PID != tc.expectedPID {
					t.Errorf("First process PID = %d, want %d", processes[0].PID, tc.expectedPID)
				}
				if processes[0].Name != tc.expectedName {
					t.Errorf("First process Name = %q, want %q", processes[0].Name, tc.expectedName)
				}
			}
		})
	}
}

// TestIntegration_ParseCSVLine 測試 CSV 行解析
func TestIntegration_ParseCSVLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected []string
	}{
		{
			name:     "Standard CSV",
			line:     `"Kiro.exe","12345","Console","1","123,456 K"`,
			expected: []string{"Kiro.exe", "12345", "Console", "1", "123,456 K"},
		},
		{
			name:     "Empty line",
			line:     "",
			expected: []string{},
		},
		{
			name:     "Single field",
			line:     `"test"`,
			expected: []string{"test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseCSVLine(tc.line)

			if len(result) != len(tc.expected) {
				t.Errorf("parseCSVLine() returned %d fields, want %d", len(result), len(tc.expected))
				return
			}

			for i, field := range result {
				if field != tc.expected[i] {
					t.Errorf("Field %d = %q, want %q", i, field, tc.expected[i])
				}
			}
		})
	}
}

// TestIntegration_PublicAPIWithShield 測試公開 API 使用 Shield 正常運作
// 驗證 GetKiroProcesses 和 IsKiroRunning 能正常執行
// **Validates: Requirements 4.1**
func TestIntegration_PublicAPIWithShield(t *testing.T) {
	shield.Init()

	// 測試 GetKiroProcesses
	processes, err := GetKiroProcesses()
	if err != nil {
		t.Fatalf("GetKiroProcesses() returned error: %v", err)
	}
	t.Logf("GetKiroProcesses() returned %d processes", len(processes))

	// 測試 IsKiroRunning
	isRunning := IsKiroRunning()
	t.Logf("IsKiroRunning() = %v", isRunning)

	// 測試 GetKiroProcessCount
	count := GetKiroProcessCount()
	t.Logf("GetKiroProcessCount() = %d", count)

	// 驗證一致性
	if isRunning != (count > 0) {
		t.Errorf("IsKiroRunning() = %v but GetKiroProcessCount() = %d", isRunning, count)
	}
	if len(processes) != count {
		t.Errorf("GetKiroProcesses() returned %d but GetKiroProcessCount() = %d", len(processes), count)
	}
}
