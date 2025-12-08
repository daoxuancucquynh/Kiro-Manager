package shield

import (
	"math/rand"
	"os/exec"
	"reflect"
	"testing"
	"testing/quick"
)

// **Feature: shield-protection, Property 3: CmdBuilder Equivalence**
// *For any* command name and arguments, CmdBuilder.Build() SHALL produce a command
// with Path and Args equivalent to exec.Command() with the decoded values.
// **Validates: Requirements 2.1, 2.3**
func TestProperty_CmdBuilderEquivalence(t *testing.T) {
	builder := NewCmdBuilder()
	codec := NewCodec()

	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// 生成隨機命令名稱（非空）
		nameLen := r.Intn(20) + 1
		name := generateCmdName(r, nameLen)

		// 生成隨機參數（0-5 個）
		argCount := r.Intn(6)
		args := make([]string, argCount)
		for i := range args {
			argLen := r.Intn(30) + 1
			args[i] = generateCmdArg(r, argLen)
		}

		// 編碼命令名稱和參數
		encodedName := codec.Encode(name)
		encodedArgs := make([][]byte, len(args))
		for i, arg := range args {
			encodedArgs[i] = codec.Encode(arg)
		}

		// 使用 CmdBuilder 建構命令
		builtCmd := builder.Build(encodedName, encodedArgs...)
		if builtCmd == nil {
			t.Logf("Build returned nil for name=%q", name)
			return false
		}

		// 使用 exec.Command 直接建構命令
		directCmd := exec.Command(name, args...)

		// Property 3: Build(Encode(name), Encode(args)...).Path = exec.Command(name, args...).Path
		if builtCmd.Path != directCmd.Path {
			t.Logf("Path mismatch: built=%q, direct=%q", builtCmd.Path, directCmd.Path)
			return false
		}

		// Property 3: Build(Encode(name), Encode(args)...).Args = exec.Command(name, args...).Args
		if !reflect.DeepEqual(builtCmd.Args, directCmd.Args) {
			t.Logf("Args mismatch: built=%v, direct=%v", builtCmd.Args, directCmd.Args)
			return false
		}

		return true
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// TestProperty_CmdBuilderWithRawArgsEquivalence 測試 BuildWithRawArgs 的等價性
func TestProperty_CmdBuilderWithRawArgsEquivalence(t *testing.T) {
	builder := NewCmdBuilder()
	codec := NewCodec()

	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// 生成隨機命令名稱（非空）
		nameLen := r.Intn(20) + 1
		name := generateCmdName(r, nameLen)

		// 生成隨機原始參數（0-5 個）
		argCount := r.Intn(6)
		rawArgs := make([]string, argCount)
		for i := range rawArgs {
			argLen := r.Intn(30) + 1
			rawArgs[i] = generateCmdArg(r, argLen)
		}

		// 編碼命令名稱
		encodedName := codec.Encode(name)

		// 使用 CmdBuilder.BuildWithRawArgs 建構命令
		builtCmd := builder.BuildWithRawArgs(encodedName, rawArgs...)
		if builtCmd == nil {
			t.Logf("BuildWithRawArgs returned nil for name=%q", name)
			return false
		}

		// 使用 exec.Command 直接建構命令
		directCmd := exec.Command(name, rawArgs...)

		// 驗證 Path 相等
		if builtCmd.Path != directCmd.Path {
			t.Logf("Path mismatch: built=%q, direct=%q", builtCmd.Path, directCmd.Path)
			return false
		}

		// 驗證 Args 相等
		if !reflect.DeepEqual(builtCmd.Args, directCmd.Args) {
			t.Logf("Args mismatch: built=%v, direct=%v", builtCmd.Args, directCmd.Args)
			return false
		}

		return true
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// TestProperty_CmdBuilderWithKnownCommands 測試已知命令的等價性
func TestProperty_CmdBuilderWithKnownCommands(t *testing.T) {
	builder := NewCmdBuilder()
	codec := NewCodec()

	// 測試實際使用的命令
	testCases := []struct {
		name string
		args []string
	}{
		{"reg", []string{"query", "HKLM\\SOFTWARE\\Microsoft\\Cryptography", "/v", "MachineGuid"}},
		{"tasklist", []string{"/FI", "IMAGENAME eq Kiro.exe", "/FO", "CSV", "/NH"}},
		{"taskkill", []string{"/PID", "12345", "/F"}},
	}

	for _, tc := range testCases {
		// 編碼
		encodedName := codec.Encode(tc.name)
		encodedArgs := make([][]byte, len(tc.args))
		for i, arg := range tc.args {
			encodedArgs[i] = codec.Encode(arg)
		}

		// 使用 CmdBuilder 建構
		builtCmd := builder.Build(encodedName, encodedArgs...)
		if builtCmd == nil {
			t.Errorf("Build returned nil for %q", tc.name)
			continue
		}

		// 使用 exec.Command 直接建構
		directCmd := exec.Command(tc.name, tc.args...)

		// 驗證等價性
		if builtCmd.Path != directCmd.Path {
			t.Errorf("Path mismatch for %q: built=%q, direct=%q",
				tc.name, builtCmd.Path, directCmd.Path)
		}

		if !reflect.DeepEqual(builtCmd.Args, directCmd.Args) {
			t.Errorf("Args mismatch for %q: built=%v, direct=%v",
				tc.name, builtCmd.Args, directCmd.Args)
		}
	}
}

// TestCmdBuilder_EmptyName 測試空命令名稱邊界情況
func TestCmdBuilder_EmptyName(t *testing.T) {
	builder := NewCmdBuilder()

	// 空 slice 應該返回 nil
	cmd := builder.Build([]byte{})
	if cmd != nil {
		t.Error("Expected nil for empty name")
	}

	// nil 應該返回 nil
	cmd = builder.Build(nil)
	if cmd != nil {
		t.Error("Expected nil for nil name")
	}

	// BuildWithRawArgs 也應該返回 nil
	cmd = builder.BuildWithRawArgs([]byte{}, "arg1")
	if cmd != nil {
		t.Error("Expected nil for empty name in BuildWithRawArgs")
	}
}

// TestCmdBuilder_NoArgs 測試無參數的命令
func TestCmdBuilder_NoArgs(t *testing.T) {
	builder := NewCmdBuilder()
	codec := NewCodec()

	name := "echo"
	encodedName := codec.Encode(name)

	builtCmd := builder.Build(encodedName)
	directCmd := exec.Command(name)

	if builtCmd == nil {
		t.Fatal("Build returned nil")
	}

	if builtCmd.Path != directCmd.Path {
		t.Errorf("Path mismatch: built=%q, direct=%q", builtCmd.Path, directCmd.Path)
	}

	if !reflect.DeepEqual(builtCmd.Args, directCmd.Args) {
		t.Errorf("Args mismatch: built=%v, direct=%v", builtCmd.Args, directCmd.Args)
	}
}

// generateCmdName 生成有效的命令名稱
func generateCmdName(r *rand.Rand, length int) string {
	// 命令名稱通常只包含字母、數字和一些特殊字元
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// generateCmdArg 生成命令參數
func generateCmdArg(r *rand.Rand, length int) string {
	// 參數可以包含更多字元，包括路徑分隔符、空格等
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/\\-_.= "
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
