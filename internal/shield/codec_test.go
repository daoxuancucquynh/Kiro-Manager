package shield

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// **Feature: shield-protection, Property 1: Codec Round-Trip**
// *For any* string input, encoding and then decoding with StringCodec SHALL
// produce a result identical to the original input.
// **Validates: Requirements 1.3**
func TestProperty_CodecRoundTrip(t *testing.T) {
	codec := NewCodec()

	f := func(s string) bool {
		// Property 1: Decode(Encode(s)) = s
		encoded := codec.Encode(s)
		decoded := codec.Decode(encoded)
		return decoded == s
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// **Feature: shield-protection, Property 2: Encoding Produces Different Output**
// *For any* non-empty string input, the encoded byte slice SHALL NOT contain
// the original string bytes in the same sequence.
// **Validates: Requirements 1.1**
func TestProperty_EncodingProducesDifferentOutput(t *testing.T) {
	codec := NewCodec()

	f := func(s string) bool {
		// 跳過空字串（空字串編碼後仍為空）
		if s == "" {
			return true
		}

		encoded := codec.Encode(s)
		original := []byte(s)

		// Property 2: Encode(s) ≠ []byte(s)
		// 編碼後的結果不應該與原始 bytes 相同
		return !bytes.Equal(encoded, original)
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// TestProperty_EncodingProducesDifferentOutput_WithGenerator 使用自定義生成器
// 確保測試涵蓋各種字串類型
func TestProperty_EncodingProducesDifferentOutput_WithGenerator(t *testing.T) {
	codec := NewCodec()

	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// 生成非空隨機字串
		length := r.Intn(100) + 1 // 1-100 字元
		s := generateRandomString(r, length)

		encoded := codec.Encode(s)
		original := []byte(s)

		// Property 2: Encode(s) ≠ []byte(s)
		return !bytes.Equal(encoded, original)
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// TestProperty_CodecRoundTrip_WithSpecialChars 測試特殊字元的 round-trip
func TestProperty_CodecRoundTrip_WithSpecialChars(t *testing.T) {
	codec := NewCodec()

	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// 生成包含特殊字元的字串
		s := generateStringWithSpecialChars(r)

		// Property 1: Decode(Encode(s)) = s
		encoded := codec.Encode(s)
		decoded := codec.Decode(encoded)
		return decoded == s
	}

	config := &quick.Config{
		MaxCount: 100,
		Values: func(values []reflect.Value, rand *rand.Rand) {
			values[0] = reflect.ValueOf(rand.Int63())
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}

// generateRandomString 生成隨機字串
func generateRandomString(r *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// generateStringWithSpecialChars 生成包含特殊字元的字串
func generateStringWithSpecialChars(r *rand.Rand) string {
	// 包含各種特殊字元：路徑分隔符、空格、中文等
	specialChars := []string{
		"\\", "/", " ", "\t", "\n",
		"HKLM\\SOFTWARE\\Microsoft\\Cryptography",
		"MachineGuid",
		"tasklist",
		"taskkill",
		"/FI",
		"IMAGENAME eq Kiro.exe",
		"中文測試",
		"日本語テスト",
		"emoji: 🎉",
	}

	// 隨機選擇一個特殊字串或組合
	if r.Intn(2) == 0 {
		return specialChars[r.Intn(len(specialChars))]
	}

	// 組合多個特殊字串
	result := ""
	count := r.Intn(3) + 1
	for i := 0; i < count; i++ {
		result += specialChars[r.Intn(len(specialChars))]
	}
	return result
}

// TestCodec_EmptyString 測試空字串邊界情況
func TestCodec_EmptyString(t *testing.T) {
	codec := NewCodec()

	// 空字串編碼
	encoded := codec.Encode("")
	if len(encoded) != 0 {
		t.Errorf("Expected empty slice for empty string, got %v", encoded)
	}

	// 空 slice 解碼
	decoded := codec.Decode([]byte{})
	if decoded != "" {
		t.Errorf("Expected empty string for empty slice, got %q", decoded)
	}

	// nil slice 解碼
	decoded = codec.Decode(nil)
	if decoded != "" {
		t.Errorf("Expected empty string for nil slice, got %q", decoded)
	}
}

// TestCodec_KnownValues 測試已知值（確保 XOR key 正確）
func TestCodec_KnownValues(t *testing.T) {
	codec := NewCodec()

	testCases := []struct {
		input    string
		expected []byte
	}{
		// 'r' = 0x72, 0x72 ^ 0x5A = 0x28
		// 'e' = 0x65, 0x65 ^ 0x5A = 0x3F
		// 'g' = 0x67, 0x67 ^ 0x5A = 0x3D
		{"reg", []byte{0x28, 0x3F, 0x3D}},
	}

	for _, tc := range testCases {
		encoded := codec.Encode(tc.input)
		if !bytes.Equal(encoded, tc.expected) {
			t.Errorf("Encode(%q) = %v, expected %v", tc.input, encoded, tc.expected)
		}

		// 驗證 round-trip
		decoded := codec.Decode(encoded)
		if decoded != tc.input {
			t.Errorf("Decode(Encode(%q)) = %q, expected %q", tc.input, decoded, tc.input)
		}
	}
}
