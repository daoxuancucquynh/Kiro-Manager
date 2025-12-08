package shield

// XOR 編碼的 key
const xorKey byte = 0x5A

// StringCodec 使用 XOR 編碼敏感字串
type StringCodec struct {
	key byte
}

// NewCodec 建立新的編碼器，使用預設 key (0x5A)
func NewCodec() *StringCodec {
	return &StringCodec{key: xorKey}
}

// Encode 編碼字串為 byte slice
// 對每個 byte 進行 XOR 運算
func (c *StringCodec) Encode(s string) []byte {
	if s == "" {
		return []byte{}
	}

	data := []byte(s)
	encoded := make([]byte, len(data))
	for i, b := range data {
		encoded[i] = b ^ c.key
	}
	return encoded
}

// Decode 解碼 byte slice 為字串
// XOR 是對稱運算，解碼使用相同的 key
func (c *StringCodec) Decode(encoded []byte) string {
	if len(encoded) == 0 {
		return ""
	}

	decoded := make([]byte, len(encoded))
	for i, b := range encoded {
		decoded[i] = b ^ c.key
	}
	return string(decoded)
}
