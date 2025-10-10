package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorStruct(t *testing.T) {
	// 测试Error结构体
	err := Error{
		Code:    404,
		Message: "not found",
	}

	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "not found", err.Message)
	assert.Equal(t, 404, err.StatusCode())
	assert.Equal(t, "[404] not found", err.Error())
}

func TestErrorMethods(t *testing.T) {
	// 测试Error的方法
	err := Error{
		Code:    500,
		Message: "internal server error",
	}

	// 测试StatusCode方法
	assert.Equal(t, 500, err.StatusCode())

	// 测试Error方法
	expectedErrorString := "[500] internal server error"
	assert.Equal(t, expectedErrorString, err.Error())
}

func TestPredefinedErrors(t *testing.T) {
	// 测试预定义的错误
	assert.Equal(t, http.StatusUnauthorized, ErrUnauthorized.Code)
	assert.Equal(t, "unauthorized", ErrUnauthorized.Message)

	assert.Equal(t, http.StatusNotFound, ErrAttachmentNotExist.Code)
	assert.Equal(t, "attachment not exist", ErrAttachmentNotExist.Message)

	assert.Equal(t, http.StatusForbidden, ErrNotAttachmentOwner.Code)
	assert.Equal(t, "not attachment owner", ErrNotAttachmentOwner.Message)
}

func TestAuthenticationErrors(t *testing.T) {
	// 测试身份认证相关错误
	assert.Equal(t, "额度不足", ErrQuotaExceeded.Error())
	assert.Equal(t, "调用语言模型失败", ErrLLMCallFailed.Error())
	assert.Equal(t, "empty password", ErrEmptyPassword.Error())
	assert.Equal(t, "empty email", ErrEmptyEmail.Error())
	assert.Equal(t, "same email", ErrSameEmail.Error())
	assert.Equal(t, "email exists, please use another email", ErrEmailExists.Error())
	assert.Equal(t, "user not exists", ErrUserNotExists.Error())
	assert.Equal(t, "forbidden access", ErrForbidden.Error())
	assert.Equal(t, "user not allow login", ErrUserNotAllowLogin.Error())
	assert.Equal(t, "user not allow signup", ErrUserNotAllowSignup.Error())
	assert.Equal(t, "user not activated", ErrNotActivated.Error())
	assert.Equal(t, "token required", ErrTokenRequired.Error())
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "bad token", ErrBadToken.Error())
	assert.Equal(t, "token expired", ErrTokenExpired.Error())
	assert.Equal(t, "email required", ErrEmailRequired.Error())
}

func TestResourceErrors(t *testing.T) {
	// 测试资源/数据处理相关错误
	assert.Equal(t, "not found", ErrNotFound.Error())
	assert.Equal(t, "not changed", ErrNotChanged.Error())
	assert.Equal(t, "with invalid view", ErrInvalidView.Error())
}

func TestPermissionErrors(t *testing.T) {
	// 测试权限与逻辑控制相关错误
	assert.Equal(t, "only super user can do this", ErrOnlySuperUser.Error())
	assert.Equal(t, "invalid primary key", ErrInvalidPrimaryKey.Error())
}

func TestToolErrors(t *testing.T) {
	// 测试工具相关错误
	assert.Equal(t, "invalid tool list response format", ErrInvalidToolListFormat.Error())
	assert.Equal(t, "invalid tool format", ErrInvalidToolFormat.Error())
	assert.Equal(t, "tool not found", ErrToolNotFound.Error())
	assert.Equal(t, "invalid tool parameters", ErrInvalidToolParams.Error())
}

func TestJSONRPCErrors(t *testing.T) {
	// 测试JSON-RPC相关错误
	assert.Equal(t, "failed to parse JSON-RPC message", ErrParseJSONRPC.Error())
	assert.Equal(t, "invalid JSON-RPC format", ErrInvalidJSONRPCFormat.Error())
	assert.Equal(t, "invalid JSON-RPC response", ErrInvalidJSONRPCResponse.Error())
	assert.Equal(t, "invalid JSON-RPC request", ErrInvalidJSONRPCRequest.Error())
	assert.Equal(t, "invalid JSON-RPC parameters", ErrInvalidJSONRPCParams.Error())
}

func TestResourceManagerErrors(t *testing.T) {
	// 测试资源管理器错误
	assert.Equal(t, "invalid resource format", ErrInvalidResourceFormat.Error())
	assert.Equal(t, "resource not found", ErrResourceNotFound.Error())
}

func TestPromptErrors(t *testing.T) {
	// 测试提示相关错误
	assert.Equal(t, "invalid prompt format", ErrInvalidPromptFormat.Error())
	assert.Equal(t, "prompt not found", ErrPromptNotFound.Error())
}

func TestToolManagerErrors(t *testing.T) {
	// 测试工具管理器错误
	assert.Equal(t, "tool name cannot be empty", ErrEmptyToolName.Error())
	assert.Equal(t, "tool already registered", ErrToolAlreadyRegistered.Error())
	assert.Equal(t, "tool execution failed", ErrToolExecutionFailed.Error())
}

func TestResourceManagerSpecificErrors(t *testing.T) {
	// 测试资源管理器特定错误
	assert.Equal(t, "resource URI cannot be empty", ErrEmptyResourceURI.Error())
}

func TestPromptManagerErrors(t *testing.T) {
	// 测试提示管理器错误
	assert.Equal(t, "prompt name cannot be empty", ErrEmptyPromptName.Error())
}

func TestLifecycleManagerErrors(t *testing.T) {
	// 测试生命周期管理器错误
	assert.Equal(t, "session already initialized", ErrSessionAlreadyInitialized.Error())
	assert.Equal(t, "session not initialized", ErrSessionNotInitialized.Error())
}

func TestParameterErrors(t *testing.T) {
	// 测试参数错误
	assert.Equal(t, "invalid parameters", ErrInvalidParams.Error())
	assert.Equal(t, "missing required parameters", ErrMissingParams.Error())
}

func TestClientErrors(t *testing.T) {
	// 测试客户端错误
	assert.Equal(t, "client already initialized", ErrAlreadyInitialized.Error())
	assert.Equal(t, "client not initialized", ErrNotInitialized.Error())
	assert.Equal(t, "invalid server URL", ErrInvalidServerURL.Error())
}

func TestErrorTypes(t *testing.T) {
	// 测试错误类型
	assert.IsType(t, &Error{}, ErrUnauthorized)
	assert.IsType(t, &Error{}, ErrAttachmentNotExist)
	assert.IsType(t, &Error{}, ErrNotAttachmentOwner)

	// 测试普通error类型
	assert.IsType(t, (*error)(nil), &ErrQuotaExceeded)
	assert.IsType(t, (*error)(nil), &ErrLLMCallFailed)
	assert.IsType(t, (*error)(nil), &ErrEmptyPassword)
}

func TestErrorConsistency(t *testing.T) {
	// 测试错误消息的一致性
	// 所有错误消息都不应该为空
	allErrors := []error{
		ErrQuotaExceeded, ErrLLMCallFailed, ErrEmptyPassword, ErrEmptyEmail,
		ErrSameEmail, ErrEmailExists, ErrUserNotExists, ErrForbidden,
		ErrUserNotAllowLogin, ErrUserNotAllowSignup, ErrNotActivated,
		ErrTokenRequired, ErrInvalidToken, ErrBadToken, ErrTokenExpired,
		ErrEmailRequired, ErrNotFound, ErrNotChanged, ErrInvalidView,
		ErrOnlySuperUser, ErrInvalidPrimaryKey, ErrInvalidToolListFormat,
		ErrInvalidToolFormat, ErrToolNotFound, ErrInvalidToolParams,
		ErrParseJSONRPC, ErrInvalidJSONRPCFormat, ErrInvalidJSONRPCResponse,
		ErrInvalidJSONRPCRequest, ErrInvalidJSONRPCParams, ErrInvalidResourceFormat,
		ErrResourceNotFound, ErrInvalidPromptFormat, ErrPromptNotFound,
		ErrEmptyToolName, ErrToolAlreadyRegistered, ErrToolExecutionFailed,
		ErrEmptyResourceURI, ErrEmptyPromptName, ErrSessionAlreadyInitialized,
		ErrSessionNotInitialized, ErrInvalidParams, ErrMissingParams,
		ErrAlreadyInitialized, ErrNotInitialized, ErrInvalidServerURL,
	}

	for _, err := range allErrors {
		assert.NotEmpty(t, err.Error(), "错误消息不应该为空: %T", err)
	}
}

func TestErrorUniqueness(t *testing.T) {
	// 测试错误消息的唯一性
	allErrorMessages := []string{
		ErrQuotaExceeded.Error(), ErrLLMCallFailed.Error(), ErrEmptyPassword.Error(),
		ErrEmptyEmail.Error(), ErrSameEmail.Error(), ErrEmailExists.Error(),
		ErrUserNotExists.Error(), ErrForbidden.Error(), ErrUserNotAllowLogin.Error(),
		ErrUserNotAllowSignup.Error(), ErrNotActivated.Error(), ErrTokenRequired.Error(),
		ErrInvalidToken.Error(), ErrBadToken.Error(), ErrTokenExpired.Error(),
		ErrEmailRequired.Error(), ErrNotFound.Error(), ErrNotChanged.Error(),
		ErrInvalidView.Error(), ErrOnlySuperUser.Error(), ErrInvalidPrimaryKey.Error(),
		ErrInvalidToolListFormat.Error(), ErrInvalidToolFormat.Error(), ErrToolNotFound.Error(),
		ErrInvalidToolParams.Error(), ErrParseJSONRPC.Error(), ErrInvalidJSONRPCFormat.Error(),
		ErrInvalidJSONRPCResponse.Error(), ErrInvalidJSONRPCRequest.Error(), ErrInvalidJSONRPCParams.Error(),
		ErrInvalidResourceFormat.Error(), ErrResourceNotFound.Error(), ErrInvalidPromptFormat.Error(),
		ErrPromptNotFound.Error(), ErrEmptyToolName.Error(), ErrToolAlreadyRegistered.Error(),
		ErrToolExecutionFailed.Error(), ErrEmptyResourceURI.Error(), ErrEmptyPromptName.Error(),
		ErrSessionAlreadyInitialized.Error(), ErrSessionNotInitialized.Error(), ErrInvalidParams.Error(),
		ErrMissingParams.Error(), ErrAlreadyInitialized.Error(), ErrNotInitialized.Error(),
		ErrInvalidServerURL.Error(),
	}

	// 检查是否有重复的错误消息
	seen := make(map[string]bool)
	for _, message := range allErrorMessages {
		assert.False(t, seen[message], "错误消息重复: %s", message)
		seen[message] = true
	}
}

func TestErrorStatusCode(t *testing.T) {
	// 测试Error结构体的StatusCode方法
	testCases := []struct {
		code     int
		message  string
		expected int
	}{
		{200, "ok", 200},
		{400, "bad request", 400},
		{401, "unauthorized", 401},
		{403, "forbidden", 403},
		{404, "not found", 404},
		{500, "internal server error", 500},
	}

	for _, tc := range testCases {
		err := Error{Code: tc.code, Message: tc.message}
		assert.Equal(t, tc.expected, err.StatusCode())
	}
}

func TestErrorString(t *testing.T) {
	// 测试Error结构体的String方法
	testCases := []struct {
		code     int
		message  string
		expected string
	}{
		{200, "ok", "[200] ok"},
		{400, "bad request", "[400] bad request"},
		{401, "unauthorized", "[401] unauthorized"},
		{403, "forbidden", "[403] forbidden"},
		{404, "not found", "[404] not found"},
		{500, "internal server error", "[500] internal server error"},
	}

	for _, tc := range testCases {
		err := Error{Code: tc.code, Message: tc.message}
		assert.Equal(t, tc.expected, err.Error())
	}
}
