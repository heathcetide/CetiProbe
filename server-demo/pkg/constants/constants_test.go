package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentConstants(t *testing.T) {
	// 测试配置缓存相关常量
	assert.Equal(t, "CONFIG_CACHE_SIZE", ENV_CONFIG_CACHE_SIZE)
	assert.Equal(t, "CONFIG_CACHE_EXPIRED", ENV_CONFIG_CACHE_EXPIRED)

	// 测试会话相关常量
	assert.Equal(t, "SESSION_FIELD", ENV_SESSION_FIELD)
	assert.Equal(t, "SESSION_SECRET", ENV_SESSION_SECRET)
	assert.Equal(t, "SESSION_EXPIRE_DAYS", ENV_SESSION_EXPIRE_DAYS)

	// 测试数据库相关常量
	assert.Equal(t, "DB_DRIVER", ENV_DB_DRIVER)
	assert.Equal(t, "DSN", ENV_DSN)

	// 测试字段常量
	assert.Equal(t, "_lingyu_db", DbField)
	assert.Equal(t, "_lingyu_uid", UserField)
	assert.Equal(t, "_lingyu_gid", GroupField)
	assert.Equal(t, "_lingyu_tz", TzField)
	assert.Equal(t, "_lingyu_assets", AssetsField)
	assert.Equal(t, "_lingyu_templates", TemplatesField)
}

func TestKeyConstants(t *testing.T) {
	// 测试过期时间相关常量
	assert.Equal(t, "VERIFY_EMAIL_EXPIRED", KEY_VERIFY_EMAIL_EXPIRED)
	assert.Equal(t, "AUTH_TOKEN_EXPIRED", KEY_AUTH_TOKEN_EXPIRED)

	// 测试站点信息相关常量
	assert.Equal(t, "SITE_NAME", KEY_SITE_NAME)
	assert.Equal(t, "SITE_ADMIN", KEY_SITE_ADMIN)
	assert.Equal(t, "SITE_URL", KEY_SITE_URL)
	assert.Equal(t, "SITE_KEYWORDS", KEY_SITE_KEYWORDS)
	assert.Equal(t, "SITE_DESCRIPTION", KEY_SITE_DESCRIPTION)
	assert.Equal(t, "SITE_GA", KEY_SITE_GA)

	// 测试站点资源相关常量
	assert.Equal(t, "SITE_LOGO_URL", KEY_SITE_LOGO_URL)
	assert.Equal(t, "SITE_FAVICON_URL", KEY_SITE_FAVICON_URL)
	assert.Equal(t, "SITE_TERMS_URL", KEY_SITE_TERMS_URL)
	assert.Equal(t, "SITE_PRIVACY_URL", KEY_SITE_PRIVACY_URL)

	// 测试认证相关常量
	assert.Equal(t, "SITE_SIGNIN_URL", KEY_SITE_SIGNIN_URL)
	assert.Equal(t, "SITE_SIGNUP_URL", KEY_SITE_SIGNUP_URL)
	assert.Equal(t, "SITE_LOGOUT_URL", KEY_SITE_LOGOUT_URL)
	assert.Equal(t, "SITE_RESET_PASSWORD_URL", KEY_SITE_RESET_PASSWORD_URL)

	// 测试API相关常量
	assert.Equal(t, "SITE_SIGNIN_API", KEY_SITE_SIGNIN_API)
	assert.Equal(t, "SITE_SIGNUP_API", KEY_SITE_SIGNUP_API)
	assert.Equal(t, "SITE_RESET_PASSWORD_DONE_API", KEY_SITE_RESET_PASSWORD_DONE_API)

	// 测试其他配置常量
	assert.Equal(t, "SITE_LOGIN_NEXT", KEY_SITE_LOGIN_NEXT)
	assert.Equal(t, "SITE_USER_ID_TYPE", KEY_SITE_USER_ID_TYPE)
	assert.Equal(t, "USER_ACTIVATED", KEY_USER_ACTIVATED)
}

func TestStaticConstants(t *testing.T) {
	// 测试静态资源相关常量
	assert.Equal(t, "STATIC_PREFIX", ENV_STATIC_PREFIX)
	assert.Equal(t, "STATIC_ROOT", ENV_STATIC_ROOT)
}

func TestConstantsValues(t *testing.T) {
	// 验证常量值不为空
	assert.NotEmpty(t, ENV_CONFIG_CACHE_SIZE)
	assert.NotEmpty(t, ENV_CONFIG_CACHE_EXPIRED)
	assert.NotEmpty(t, ENV_SESSION_FIELD)
	assert.NotEmpty(t, ENV_SESSION_SECRET)
	assert.NotEmpty(t, ENV_SESSION_EXPIRE_DAYS)
	assert.NotEmpty(t, ENV_DB_DRIVER)
	assert.NotEmpty(t, ENV_DSN)
	assert.NotEmpty(t, DbField)
	assert.NotEmpty(t, UserField)
	assert.NotEmpty(t, GroupField)
	assert.NotEmpty(t, TzField)
	assert.NotEmpty(t, AssetsField)
	assert.NotEmpty(t, TemplatesField)
}

func TestConstantsFormat(t *testing.T) {
	// 验证常量格式正确性
	// 环境变量常量应该全大写
	assert.Equal(t, "CONFIG_CACHE_SIZE", ENV_CONFIG_CACHE_SIZE)
	assert.Equal(t, "CONFIG_CACHE_EXPIRED", ENV_CONFIG_CACHE_EXPIRED)
	assert.Equal(t, "SESSION_FIELD", ENV_SESSION_FIELD)
	assert.Equal(t, "SESSION_SECRET", ENV_SESSION_SECRET)
	assert.Equal(t, "SESSION_EXPIRE_DAYS", ENV_SESSION_EXPIRE_DAYS)
	assert.Equal(t, "DB_DRIVER", ENV_DB_DRIVER)
	assert.Equal(t, "DSN", ENV_DSN)
	assert.Equal(t, "STATIC_PREFIX", ENV_STATIC_PREFIX)
	assert.Equal(t, "STATIC_ROOT", ENV_STATIC_ROOT)
}

func TestConstantsUniqueness(t *testing.T) {
	// 验证常量值的唯一性
	allConstants := []string{
		ENV_CONFIG_CACHE_SIZE, ENV_CONFIG_CACHE_EXPIRED,
		ENV_SESSION_FIELD, ENV_SESSION_SECRET, ENV_SESSION_EXPIRE_DAYS,
		ENV_DB_DRIVER, ENV_DSN,
		DbField, UserField, GroupField, TzField, AssetsField, TemplatesField,
		KEY_VERIFY_EMAIL_EXPIRED, KEY_AUTH_TOKEN_EXPIRED,
		KEY_SITE_NAME, KEY_SITE_ADMIN, KEY_SITE_URL, KEY_SITE_KEYWORDS,
		KEY_SITE_DESCRIPTION, KEY_SITE_GA,
		KEY_SITE_LOGO_URL, KEY_SITE_FAVICON_URL, KEY_SITE_TERMS_URL,
		KEY_SITE_PRIVACY_URL, KEY_SITE_SIGNIN_URL, KEY_SITE_SIGNUP_URL,
		KEY_SITE_LOGOUT_URL, KEY_SITE_RESET_PASSWORD_URL,
		KEY_SITE_SIGNIN_API, KEY_SITE_SIGNUP_API, KEY_SITE_RESET_PASSWORD_DONE_API,
		KEY_SITE_LOGIN_NEXT, KEY_SITE_USER_ID_TYPE, KEY_USER_ACTIVATED,
		ENV_STATIC_PREFIX, ENV_STATIC_ROOT,
	}

	// 检查是否有重复的常量值
	seen := make(map[string]bool)
	for _, constant := range allConstants {
		assert.False(t, seen[constant], "常量值 %s 重复了", constant)
		seen[constant] = true
	}
}
