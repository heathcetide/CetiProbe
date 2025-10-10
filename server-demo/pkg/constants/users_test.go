package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSignalConstants(t *testing.T) {
	// 测试用户信号常量
	assert.Equal(t, "user.login", SigUserLogin)
	assert.Equal(t, "user.logout", SigUserLogout)
	assert.Equal(t, "user.create", SigUserCreate)
	assert.Equal(t, "user.verifyemail", SigUserVerifyEmail)
	assert.Equal(t, "user.resetpassword", SigUserResetPassword)
	assert.Equal(t, "user.changeemail", SigUserChangeEmail)
	assert.Equal(t, "user.changeemaildone", SigUserChangeEmailDone)
}

func TestUserSignalConstantsNotEmpty(t *testing.T) {
	// 验证所有用户信号常量不为空
	assert.NotEmpty(t, SigUserLogin)
	assert.NotEmpty(t, SigUserLogout)
	assert.NotEmpty(t, SigUserCreate)
	assert.NotEmpty(t, SigUserVerifyEmail)
	assert.NotEmpty(t, SigUserResetPassword)
	assert.NotEmpty(t, SigUserChangeEmail)
	assert.NotEmpty(t, SigUserChangeEmailDone)
}

func TestUserSignalConstantsFormat(t *testing.T) {
	// 验证用户信号常量格式
	// 所有信号都应该以"user."开头
	assert.True(t, len(SigUserLogin) > 5 && SigUserLogin[:5] == "user.")
	assert.True(t, len(SigUserLogout) > 5 && SigUserLogout[:5] == "user.")
	assert.True(t, len(SigUserCreate) > 5 && SigUserCreate[:5] == "user.")
	assert.True(t, len(SigUserVerifyEmail) > 5 && SigUserVerifyEmail[:5] == "user.")
	assert.True(t, len(SigUserResetPassword) > 5 && SigUserResetPassword[:5] == "user.")
	assert.True(t, len(SigUserChangeEmail) > 5 && SigUserChangeEmail[:5] == "user.")
	assert.True(t, len(SigUserChangeEmailDone) > 5 && SigUserChangeEmailDone[:5] == "user.")
}

func TestUserSignalConstantsUniqueness(t *testing.T) {
	// 验证用户信号常量的唯一性
	allSignals := []string{
		SigUserLogin, SigUserLogout, SigUserCreate, SigUserVerifyEmail,
		SigUserResetPassword, SigUserChangeEmail, SigUserChangeEmailDone,
	}

	// 检查是否有重复的信号
	seen := make(map[string]bool)
	for _, signal := range allSignals {
		assert.False(t, seen[signal], "用户信号 %s 重复了", signal)
		seen[signal] = true
	}
}

func TestUserSignalConstantsLength(t *testing.T) {
	// 验证用户信号常量的长度合理性
	assert.True(t, len(SigUserLogin) >= 10, "SigUserLogin 长度应该至少为10")
	assert.True(t, len(SigUserLogout) >= 11, "SigUserLogout 长度应该至少为11")
	assert.True(t, len(SigUserCreate) >= 11, "SigUserCreate 长度应该至少为11")
	assert.True(t, len(SigUserVerifyEmail) >= 16, "SigUserVerifyEmail 长度应该至少为16")
	assert.True(t, len(SigUserResetPassword) >= 18, "SigUserResetPassword 长度应该至少为18")
	assert.True(t, len(SigUserChangeEmail) >= 16, "SigUserChangeEmail 长度应该至少为16")
	assert.True(t, len(SigUserChangeEmailDone) >= 20, "SigUserChangeEmailDone 长度应该至少为20")
}

func TestUserSignalConstantsContent(t *testing.T) {
	// 验证用户信号常量的具体内容
	assert.Contains(t, SigUserLogin, "login")
	assert.Contains(t, SigUserLogout, "logout")
	assert.Contains(t, SigUserCreate, "create")
	assert.Contains(t, SigUserVerifyEmail, "verifyemail")
	assert.Contains(t, SigUserResetPassword, "resetpassword")
	assert.Contains(t, SigUserChangeEmail, "changeemail")
	assert.Contains(t, SigUserChangeEmailDone, "changeemaildone")
}

func TestUserSignalConstantsConsistency(t *testing.T) {
	// 验证用户信号常量的一致性
	// 所有信号都应该使用小写字母和点号分隔
	allSignals := []string{
		SigUserLogin, SigUserLogout, SigUserCreate, SigUserVerifyEmail,
		SigUserResetPassword, SigUserChangeEmail, SigUserChangeEmailDone,
	}

	for _, signal := range allSignals {
		// 检查是否只包含小写字母、数字和点号
		for _, char := range signal {
			assert.True(t,
				(char >= 'a' && char <= 'z') ||
					(char >= '0' && char <= '9') ||
					char == '.',
				"信号 %s 包含非法字符 %c", signal, char)
		}
	}
}

func TestUserSignalConstantsSemantic(t *testing.T) {
	// 验证用户信号常量的语义正确性
	// 登录相关
	assert.Contains(t, SigUserLogin, "login")
	assert.Contains(t, SigUserLogout, "logout")

	// 用户创建
	assert.Contains(t, SigUserCreate, "create")

	// 邮箱验证
	assert.Contains(t, SigUserVerifyEmail, "verify")
	assert.Contains(t, SigUserVerifyEmail, "email")

	// 密码重置
	assert.Contains(t, SigUserResetPassword, "reset")
	assert.Contains(t, SigUserResetPassword, "password")

	// 邮箱更改
	assert.Contains(t, SigUserChangeEmail, "change")
	assert.Contains(t, SigUserChangeEmail, "email")
	assert.Contains(t, SigUserChangeEmailDone, "change")
	assert.Contains(t, SigUserChangeEmailDone, "email")
	assert.Contains(t, SigUserChangeEmailDone, "done")
}
