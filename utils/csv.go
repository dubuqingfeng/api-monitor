package utils

import "bytes"

var UTF8BOM = []byte{239, 187, 191}

func HasBOM(in []byte) bool {
	return bytes.HasPrefix(in, UTF8BOM)
}

func StripBOM(in []byte) []byte {
	return bytes.TrimPrefix(in, UTF8BOM)
}