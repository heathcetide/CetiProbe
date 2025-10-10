package layer

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
)

// ErrorLayerInfo 存储错误层信息
type ErrorLayerInfo struct {
	// 基本信息
	Timestamp time.Time `json:"timestamp"` // 数据包捕获时间

	// 错误信息
	Error string `json:"error,omitempty"` // 错误描述
	Layer string `json:"layer,omitempty"` // 出错的层
	Fatal bool   `json:"fatal,omitempty"` // 是否为致命错误
	Code  int    `json:"code,omitempty"`  // 错误代码
}

// ExtractErrorLayerInfo 提取错误层信息并填充到ErrorLayerInfo结构体中
func ExtractErrorLayerInfo(errLayer gopacket.ErrorLayer) *ErrorLayerInfo {
	info := &ErrorLayerInfo{
		Timestamp: time.Now(),
	}

	if errLayer == nil {
		return info
	}

	// 获取错误信息
	err := errLayer.Error()
	info.Error = err.Error()
	info.Layer = fmt.Sprintf("%T", errLayer.LayerContents())

	// 检查是否为致命错误
	// 这里我们可以根据错误类型判断是否为致命错误
	// 通常校验和错误不是致命的，但格式错误可能是致命的
	if strings.Contains(strings.ToLower(err.Error()), "checksum") {
		info.Fatal = false
		info.Code = 1 // 校验和错误
	} else {
		info.Fatal = true
		info.Code = 2 // 其他错误
	}

	return info
}

// PrintErrorLayerDetails 打印错误层详细信息
func PrintErrorLayerDetails(errInfo *ErrorLayerInfo) {
	fmt.Println("  Error Layer 详细信息:")
	fmt.Printf("    错误: %s\n", errInfo.Error)
	fmt.Printf("    层: %s\n", errInfo.Layer)
	fmt.Printf("    类型: %s\n", map[bool]string{true: "致命错误", false: "非致命错误"}[errInfo.Fatal])
	if errInfo.Code != 0 {
		fmt.Printf("    错误代码: %d\n", errInfo.Code)
	}
}
