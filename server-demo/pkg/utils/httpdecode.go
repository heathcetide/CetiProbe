package utils

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"mime"
	"strings"
)

// DecodeBodyToText tries to decode HTTP body bytes to human-readable text based on headers.
// Supports gzip and deflate. If charset detected in Content-Type, tries to honor utf-8; other charsets are returned best-effort.
func DecodeBodyToText(body []byte, headers map[string]string) (string, string) {
	if len(body) == 0 {
		return "", ""
	}
	contentEncoding := strings.ToLower(headers["Content-Encoding"])
	var raw []byte = body
	switch contentEncoding {
	case "gzip":
		zr, err := gzip.NewReader(bytes.NewReader(body))
		if err == nil {
			out, _ := io.ReadAll(zr)
			_ = zr.Close()
			if len(out) > 0 {
				raw = out
			}
		}
	case "deflate":
		zr, err := zlib.NewReader(bytes.NewReader(body))
		if err == nil {
			out, _ := io.ReadAll(zr)
			_ = zr.Close()
			if len(out) > 0 {
				raw = out
			}
		}
	}

	ct := headers["Content-Type"]
	// Try to detect if it's likely text
	if isLikelyText(raw) {
		// Best effort: do not convert charset; Go runtime lacks generic iconv.
		return string(raw), ct
	}
	// Not text, return empty text with content-type hint
	return "", ct
}

func isLikelyText(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	printable := 0
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printable++
		}
	}
	return float64(printable)/float64(len(data)) > 0.6
}

// GetMimeType parses Content-Type and returns mime type without charset
func GetMimeType(contentType string) string {
	if contentType == "" {
		return ""
	}
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mt
}
