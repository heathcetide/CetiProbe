package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// SecureQueryBuilder 安全查询构建器
type SecureQueryBuilder struct {
	Db *gorm.DB
}

// NewSecureQueryBuilder 创建安全查询构建器
func NewSecureQueryBuilder(db *gorm.DB) *SecureQueryBuilder {
	return &SecureQueryBuilder{Db: db}
}

// SafeWhere 安全的WHERE条件构建
func (sqb *SecureQueryBuilder) SafeWhere(column string, operator string, value interface{}) *gorm.DB {
	// 验证列名
	if !isValidColumnName(column) {
		panic(fmt.Sprintf("invalid column name: %s", column))
	}

	// 验证操作符
	if !isValidOperator(operator) {
		panic(fmt.Sprintf("invalid operator: %s", operator))
	}

	// 根据操作符构建查询
	switch strings.ToUpper(operator) {
	case "=", "!=", "<", ">", "<=", ">=":
		return sqb.Db.Where(fmt.Sprintf("%s %s ?", column, operator), value)
	case "<>":
		return sqb.Db.Where(fmt.Sprintf("%s <> ?", column), value)
	case "LIKE":
		if str, ok := value.(string); ok {
			return sqb.Db.Where(fmt.Sprintf("%s LIKE ?", column), "%"+str+"%")
		}
		return sqb.Db.Where(fmt.Sprintf("%s LIKE ?", column), value)
	case "NOT LIKE":
		if str, ok := value.(string); ok {
			return sqb.Db.Where(fmt.Sprintf("%s NOT LIKE ?", column), "%"+str+"%")
		}
		return sqb.Db.Where(fmt.Sprintf("%s NOT LIKE ?", column), value)
	case "IN":
		return sqb.Db.Where(fmt.Sprintf("%s IN ?", column), value)
	case "NOT IN":
		return sqb.Db.Where(fmt.Sprintf("%s NOT IN ?", column), value)
	case "BETWEEN":
		if values, ok := value.([]interface{}); ok && len(values) == 2 {
			return sqb.Db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), values[0], values[1])
		}
		panic("BETWEEN operator requires exactly 2 values")
	case "IS NULL":
		return sqb.Db.Where(fmt.Sprintf("%s IS NULL", column))
	case "IS NOT NULL":
		return sqb.Db.Where(fmt.Sprintf("%s IS NOT NULL", column))
	default:
		panic(fmt.Sprintf("unsupported operator: %s", operator))
	}
}

// SafeOrder 安全的ORDER BY构建
func (sqb *SecureQueryBuilder) SafeOrder(column string, direction string) *gorm.DB {
	// 验证列名
	if !isValidColumnName(column) {
		panic(fmt.Sprintf("invalid column name: %s", column))
	}

	// 验证排序方向
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC" // 默认升序
	}

	return sqb.Db.Order(fmt.Sprintf("%s %s", column, direction))
}

// SafeSelect 安全的SELECT字段构建
func (sqb *SecureQueryBuilder) SafeSelect(columns []string) *gorm.DB {
	// 验证所有列名
	for _, column := range columns {
		if !isValidColumnName(column) {
			panic(fmt.Sprintf("invalid column name: %s", column))
		}
	}

	return sqb.Db.Select(columns)
}

// SafeGroup 安全的GROUP BY构建
func (sqb *SecureQueryBuilder) SafeGroup(columns []string) *gorm.DB {
	// 验证所有列名
	for _, column := range columns {
		if !isValidColumnName(column) {
			panic(fmt.Sprintf("invalid column name: %s", column))
		}
	}

	return sqb.Db.Group(strings.Join(columns, ", "))
}

// SafeHaving 安全的HAVING条件构建
func (sqb *SecureQueryBuilder) SafeHaving(condition string, args ...interface{}) *gorm.DB {
	// 验证HAVING条件中的列名
	if !isValidHavingCondition(condition) {
		panic(fmt.Sprintf("invalid HAVING condition: %s", condition))
	}

	return sqb.Db.Having(condition, args...)
}

// isValidColumnName 验证列名是否安全
func isValidColumnName(column string) bool {
	// 列名只能包含字母、数字、下划线和点
	pattern := `^[a-zA-Z_][a-zA-Z0-9_.]*$`
	matched, _ := regexp.MatchString(pattern, column)
	return matched && len(column) <= 64
}

// isValidOperator 验证操作符是否安全
func isValidOperator(operator string) bool {
	validOperators := map[string]bool{
		"=":           true,
		"!=":          true,
		"<>":          true,
		"<":           true,
		">":           true,
		"<=":          true,
		">=":          true,
		"LIKE":        true,
		"NOT LIKE":    true,
		"IN":          true,
		"NOT IN":      true,
		"BETWEEN":     true,
		"IS NULL":     true,
		"IS NOT NULL": true,
	}

	return validOperators[strings.ToUpper(operator)]
}

// isValidHavingCondition 验证HAVING条件是否安全
func isValidHavingCondition(condition string) bool {
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "UNION", "SCRIPT", "JAVASCRIPT",
	}
	upper := strings.ToUpper(condition)
	for _, kw := range dangerousKeywords {
		if strings.Contains(upper, kw) {
			return false
		}
	}
	// 允许参数占位 ?、百分号 %、单双引号和常见算术符号
	pattern := `^[a-zA-Z0-9_.,()\s=<>!%?\+\-\*/'"]+$`
	matched, _ := regexp.MatchString(pattern, condition)
	return matched
}

// SafeQuery 执行安全的原生查询
func (sqb *SecureQueryBuilder) SafeQuery(query string, args ...interface{}) *gorm.DB {
	// 验证查询语句
	if !isValidQuery(query) {
		panic(fmt.Sprintf("invalid query: %s", query))
	}

	return sqb.Db.Raw(query, args...)
}

// isValidQuery 验证查询语句是否安全
func isValidQuery(query string) bool {
	// 转换为大写进行检查
	upperQuery := strings.ToUpper(query)

	// 检查是否包含危险的关键字
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "UNION", "SCRIPT", "JAVASCRIPT",
		"TRUNCATE", "GRANT", "REVOKE", "SHUTDOWN",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperQuery, keyword) {
			return false
		}
	}

	// 检查是否以SELECT开头（只允许查询操作）
	if !strings.HasPrefix(strings.TrimSpace(upperQuery), "SELECT") {
		return false
	}

	return true
}

// SafeTransaction 安全的事务执行
func (sqb *SecureQueryBuilder) SafeTransaction(fn func(*gorm.DB) error) error {
	return sqb.Db.Transaction(func(tx *gorm.DB) error {
		// 创建新的事务查询构建器
		txBuilder := NewSecureQueryBuilder(tx)

		// 执行事务函数
		return fn(txBuilder.Db)
	})
}

// SanitizeValue 清理值，防止注入
func SanitizeValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		// 移除危险字符
		v = strings.ReplaceAll(v, "'", "''")
		v = strings.ReplaceAll(v, "\"", "\\\"")
		v = strings.ReplaceAll(v, "\\", "\\\\")
		return v
	case []string:
		// 清理字符串数组
		sanitized := make([]string, len(v))
		for i, s := range v {
			sanitized[i] = SanitizeValue(s).(string)
		}
		return sanitized
	case time.Time:
		// 时间类型直接返回
		return v
	case int, int8, int16, int32, int64:
		// 整数类型直接返回
		return v
	case uint, uint8, uint16, uint32, uint64:
		// 无符号整数类型直接返回
		return v
	case float32, float64:
		// 浮点数类型直接返回
		return v
	case bool:
		// 布尔类型直接返回
		return v
	default:
		// 其他类型转换为字符串后清理
		return SanitizeValue(fmt.Sprintf("%v", v))
	}
}

// ValidateInput 验证输入参数
func ValidateInput(input interface{}) error {
	if input == nil {
		return nil
	}
	s := fmt.Sprintf("%v", input)
	if len(s) > 10000 {
		return fmt.Errorf("input too long")
	}
	sqlPatterns := []string{
		`(?i)\bunion\s+select\b`,
		`(?i)\bdrop\s+table\b`,
		`(?i)\bdelete\s+from\b`,
		`(?i)\binsert\s+into\s+\S+`,  // insert into <tbl>
		`(?i)\bupdate\s+\S+\s+set\b`, // update <tbl> set
		`(?i)\bor\s+1\s*=\s*1\b`,
		`(?i)\band\s+1\s*=\s*1\b`,
		`(?i)\bexec\s*\(`,
	}
	for _, p := range sqlPatterns {
		if matched, _ := regexp.MatchString(p, s); matched {
			return fmt.Errorf("potentially malicious input detected")
		}
	}
	return nil
}

// SafePaginate 安全的分页查询
func (sqb *SecureQueryBuilder) SafePaginate(page, pageSize int) *gorm.DB {
	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20 // 默认每页20条
	}

	offset := (page - 1) * pageSize
	return sqb.Db.Offset(offset).Limit(pageSize)
}

// SafeCount 安全的计数查询
func (sqb *SecureQueryBuilder) SafeCount(model interface{}) (int64, error) {
	var count int64
	err := sqb.Db.Model(model).Count(&count).Error
	return count, err
}

// SafeExists 安全的存在性检查
func (sqb *SecureQueryBuilder) SafeExists(model interface{}, conditions map[string]interface{}) (bool, error) {
	query := sqb.Db.Model(model)

	// 安全地添加条件
	for column, value := range conditions {
		if !isValidColumnName(column) {
			return false, fmt.Errorf("invalid column name: %s", column)
		}
		query = query.Where(fmt.Sprintf("%s = ?", column), value)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// SafeFirst 安全的第一条记录查询
func (sqb *SecureQueryBuilder) SafeFirst(dest interface{}, conditions map[string]interface{}) error {
	query := sqb.Db.Model(dest)

	// 安全地添加条件
	for column, value := range conditions {
		if !isValidColumnName(column) {
			return fmt.Errorf("invalid column name: %s", column)
		}
		query = query.Where(fmt.Sprintf("%s = ?", column), value)
	}

	return query.First(dest).Error
}

// SafeFind 安全的批量查询
func (sqb *SecureQueryBuilder) SafeFind(dest interface{}, conditions map[string]interface{}) error {
	query := sqb.Db.Model(dest)

	// 安全地添加条件
	for column, value := range conditions {
		if !isValidColumnName(column) {
			return fmt.Errorf("invalid column name: %s", column)
		}
		query = query.Where(fmt.Sprintf("%s = ?", column), value)
	}

	return query.Find(dest).Error
}
