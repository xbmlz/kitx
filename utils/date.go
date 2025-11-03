package utils

import (
	"strings"
	"time"
)

// GetAgeByBirthday 根据生日计算年龄, 生日格式 yyyy-MM-dd
func GetAgeByBirthday(birthday string) int {
	// 检查格式是否为 yyyy-MM-dd
	if len(strings.Split(birthday, "-")) != 3 {
		return 0
	}

	// 尝试解析日期
	birthDate, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return 0
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// 如果今年的生日还没到，年龄减 1
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
