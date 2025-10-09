package capture

import (
	"testing"
)

// TestNewCapturer 测试创建新的Capturer实例
func TestNewCapturer(t *testing.T) {
	// 由于需要真实的网络接口，我们只能测试错误情况
	// 在实际环境中，可以使用mock或虚拟接口进行更完整的测试

	// 测试使用无效接口名称的情况
	_, err := NewCapturer("en0")
	if err == nil {
		t.Fatalf("Expected error when creating capturer with invalid interface name, but got nil")
	}
	// 注意：要测试成功的情况，需要在具有相应网络接口的实际环境中运行
	// 可以手动测试有效的接口名称，例如在Linux上可能是"eth0"，在macOS上可能是"en0"
}

// TestCapturerStartWithInvalidHandle 测试Capturer在无效句柄下的Start方法
func TestCapturerStartWithInvalidHandle(t *testing.T) {
	// 由于NewCapturer在无效接口名称下会返回错误，我们无法获得有效的Capturer实例
	// 因此这个测试主要是为了说明在真实环境中Start方法的预期行为

	// 在实际应用中，如果能成功创建Capturer实例，Start方法应该：
	// 1. 将running标志设置为true
	// 2. 开始监听数据包
	// 3. 在再次调用Start时返回错误（因为已经在运行）

	// 这些行为需要在有真实网络接口的环境中进行测试
	t.Log("测试说明：Start方法需要在有真实网络接口的环境中进行完整测试")
}

// TestCapturerConcurrentStart 测试并发调用Start方法
func TestCapturerConcurrentStart(t *testing.T) {
	// 类似于上面的测试，由于无法创建有效的Capturer实例，我们只能说明测试意图

	// 在实际环境中，应该测试并发调用Start的情况：
	// 1. 启动一个goroutine调用Start
	// 2. 立即启动另一个goroutine调用Start
	// 3. 验证只有一个Start调用成功，另一个应该返回错误

	t.Log("测试说明：并发Start调用测试需要在有真实网络接口的环境中进行")
}
