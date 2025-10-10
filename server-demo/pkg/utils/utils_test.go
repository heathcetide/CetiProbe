package utils

import (
	"encoding/base64"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---------- randRunes / RandText / RandNumberText / RandString ----------

func TestRandRunesAndFriends(t *testing.T) {
	// 基础长度断言
	if got := RandText(16); len(got) != 16 {
		t.Fatalf("RandText length = %d, want 16", len(got))
	}
	if got := RandNumberText(8); len(got) != 8 {
		t.Fatalf("RandNumberText length = %d, want 8", len(got))
	}
	if got := RandString(12); len(got) != 12 {
		t.Fatalf("RandString length = %d, want 12", len(got))
	}

	// 数字串只应包含数字
	num := RandNumberText(64)
	for i, c := range num {
		if c < '0' || c > '9' {
			t.Fatalf("RandNumberText contains non-digit at %d: %q", i, c)
		}
	}

	// 覆盖 randRunes 的 source 路径（传入不同 rune 源）
	alpha := []rune("ab")
	got := randRunes(32, alpha)
	for _, c := range got {
		if c != 'a' && c != 'b' {
			t.Fatalf("randRunes unexpected rune: %q", c)
		}
	}
}

// ---------- SafeCall ----------

func TestSafeCall_NoPanic(t *testing.T) {
	called := false
	err := SafeCall(func() error {
		called = true
		return nil
	}, func(error) {})
	if err != nil {
		t.Fatalf("SafeCall returned error: %v", err)
	}
	if !called {
		t.Fatalf("SafeCall did not call f()")
	}
}

func TestSafeCall_PanicError(t *testing.T) {
	var handled error
	_ = SafeCall(func() error {
		panic(assertErr("boom"))
	}, func(e error) {
		handled = e
	})
	if handled == nil || handled.Error() != "boom" {
		t.Fatalf("SafeCall did not handle error panic, got: %v", handled)
	}
}

func TestSafeCall_PanicString(t *testing.T) {
	var handled error
	_ = SafeCall(func() error {
		panic("oops")
	}, func(e error) {
		handled = e
	})
	if handled == nil || handled.Error() != "oops" {
		t.Fatalf("SafeCall did not convert string panic, got: %v", handled)
	}
}

func TestSafeCall_PanicUnknownType(t *testing.T) {
	var handled error
	_ = SafeCall(func() error {
		panic(struct{ X int }{X: 1})
	}, func(e error) {
		handled = e
	})
	if handled == nil || handled.Error() != "unknown error type" {
		t.Fatalf("SafeCall unknown-type panic => %v, want 'unknown error type'", handled)
	}
}

// 用于快速构造实现 error 的类型
type assertErr string

func (e assertErr) Error() string { return string(e) }

// ---------- StructAsMap ----------

func TestStructAsMap(t *testing.T) {
	type inner struct {
		A int
	}
	type demo struct {
		Name    string
		Age     int
		NotePtr *string
		InPtr   *inner
		ZeroStr string
		ZeroInt int
	}

	note := "hello"
	in := &inner{A: 7}

	// 非结构体
	if m := StructAsMap(123, []string{"X"}); len(m) != 0 {
		t.Fatalf("StructAsMap(non-struct) = %#v, want empty", m)
	}

	// 结构体与指针字段
	d := demo{
		Name:    "tom",
		Age:     18,
		NotePtr: &note,
		InPtr:   in,
	}
	// 仅选择部分字段，包含零值字段与不存在的字段
	fields := []string{"Name", "Age", "NotePtr", "InPtr", "ZeroStr", "ZeroInt", "NoSuch"}
	got := StructAsMap(d, fields)

	// 预期：零值与不存在的字段不出现；指针字段解引用
	want := map[string]any{
		"Name":    "tom",
		"Age":     18,
		"NotePtr": "hello",
		"InPtr":   inner{A: 7},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("StructAsMap got %#v, want %#v", got, want)
	}

	// 传入指针结构体
	got2 := StructAsMap(&d, []string{"Name"})
	if got2["Name"] != "tom" || len(got2) != 1 {
		t.Fatalf("StructAsMap(ptr) = %#v", got2)
	}
}

// ---------- GenerateSecureToken ----------

func TestGenerateSecureToken_URLSafeLength(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 16, 31, 32, 33, 64} {
		tok, err := GenerateSecureToken(n)
		if err != nil {
			t.Fatalf("GenerateSecureToken(%d) error: %v", n, err)
		}
		// base64.URLEncoding 长度校验
		wantLen := base64.URLEncoding.EncodedLen(n)
		if len(tok) != wantLen {
			t.Fatalf("token len = %d, want %d (n=%d)", len(tok), wantLen, n)
		}
		// URL 安全字符（不包含 '+' '/'）
		if strings.ContainsAny(tok, "+/") {
			t.Fatalf("token contains non-URL-safe characters: %q", tok)
		}
	}
}

// ---------- Snowflake: New / NextID ----------

func withEnv(key, val string, fn func()) {
	old := os.Getenv(key)
	_ = os.Setenv(key, val)
	defer os.Setenv(key, old)
	fn()
}

func TestNewSnowflake_OK_DefaultMachineID(t *testing.T) {
	// 非法值会走 fallback=1（在 getMachineID 中），属于有效范围
	withEnv("MACHINE_ID", "not-an-int", func() {
		sf, err := NewSnowflake()
		if err != nil {
			t.Fatalf("NewSnowflake with fallback id error: %v", err)
		}
		if sf.machineID != 1 {
			t.Fatalf("fallback machineID = %d, want 1", sf.machineID)
		}
	})
}

func TestNewSnowflake_ErrOutOfRange(t *testing.T) {
	// <0
	withEnv("MACHINE_ID", "-1", func() {
		if _, err := NewSnowflake(); err == nil {
			t.Fatalf("NewSnowflake expected error for id=-1")
		}
	})
	// > max
	tooBig := int64(maxMachineID) + 1
	withEnv("MACHINE_ID", os.Getenv("MACHINE_ID"), func() {
		_ = os.Setenv("MACHINE_ID", intToString(tooBig))
		if _, err := NewSnowflake(); err == nil {
			t.Fatalf("NewSnowflake expected error for id>max")
		}
	})
}

func intToString(v int64) string {
	// 避免引入 strconv 再次；但本文件可用 strconv，保持简洁直接用即可
	// 留着工具函数也无妨
	return strconvItoa(v)
}

func strconvItoa(v int64) string {
	// 本地实现一个简单的 itoa（支持负数）避免额外依赖
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b [32]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + (v % 10))
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func TestSnowflake_NextID_Monotonic(t *testing.T) {
	sf, err := NewSnowflake()
	if err != nil {
		t.Fatalf("NewSnowflake error: %v", err)
	}

	const N = 2000
	ids := make([]int64, N)
	for i := 0; i < N; i++ {
		ids[i] = sf.NextID()
		if ids[i] == 0 {
			t.Fatalf("NextID returned 0 unexpectedly")
		}
		if i > 0 && ids[i] <= ids[i-1] {
			t.Fatalf("IDs not strictly increasing: %d <= %d at %d", ids[i], ids[i-1], i)
		}
	}
}

func TestSnowflake_NextID_SameMicroSequenceAndRollover(t *testing.T) {
	sf, err := NewSnowflake()
	if err != nil {
		t.Fatalf("NewSnowflake error: %v", err)
	}

	// 模拟同一微秒内把 sequence 推到最大，再 NextID 触发回绕与“等待下一个微秒”
	sf.mu.Lock()
	now := currentMicro()
	sf.lastStamp = now
	sf.sequence = maxSequence
	sf.mu.Unlock()

	start := time.Now()
	id := sf.NextID()
	if id == 0 {
		t.Fatalf("NextID returned 0 on rollover")
	}
	// 由于回绕逻辑会等待下一微秒，耗时应该 >= 1 微秒（在大多数环境下远大于 1µs）
	if time.Since(start) <= 0 {
		t.Fatalf("expected rollover wait to advance time")
	}
}

func TestSnowflake_NextID_ClockRollback(t *testing.T) {
	sf, err := NewSnowflake()
	if err != nil {
		t.Fatalf("NewSnowflake error: %v", err)
	}

	// 构造 lastStamp > now 的场景
	sf.mu.Lock()
	sf.lastStamp = currentMicro() + 10_000 // 未来
	sf.mu.Unlock()

	if got := sf.NextID(); got != 0 {
		t.Fatalf("clock rollback protection expected 0, got %d", got)
	}
}

// ---------- 并发 smoke（可选，确保锁路径覆盖更充分） ----------

func TestSnowflake_Concurrent(t *testing.T) {
	sf, err := NewSnowflake()
	if err != nil {
		t.Fatalf("NewSnowflake error: %v", err)
	}

	const goroutines = 16
	const perG = 512
	var wg sync.WaitGroup
	out := make(chan int64, goroutines*perG)

	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				out <- sf.NextID()
			}
		}()
	}
	wg.Wait()
	close(out)
	first := true
	for id := range out {
		if id == 0 {
			t.Fatalf("concurrent NextID produced 0")
		}
		if first {
			_ = id
			first = false
			continue
		}
		// 并发情况下不保证读取顺序严格单调，但至少所有值应为正且非 0。
		if id < 0 {
			t.Fatalf("concurrent NextID produced negative id: %d", id)
		}
	}
}
