package middleware

import (
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// GetUserTokenVersion 获取用户Token版本号（从Redis）
func GetUserTokenVersion(rds *redis.Redis, userID int64) (int64, error) {
	key := fmt.Sprintf("user:token_version:%d", userID)
	versionStr, err := rds.Get(key)
	if err != nil {
		// 如果不存在，返回默认版本号1（兼容旧Token）
		return 1, nil
	}
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		logx.Errorf("解析Token版本号失败: %v", err)
		return 1, nil // 解析失败时返回默认版本号
	}
	return version, nil
}

// CheckTokenVersion 检查Token版本号是否有效
func CheckTokenVersion(rds *redis.Redis, userID int64, tokenVersion int64) bool {
	currentVersion, err := GetUserTokenVersion(rds, userID)
	if err != nil {
		logx.Errorf("获取用户Token版本号失败: %v", err)
		return false
	}
	// Token版本号必须 >= 当前版本号才有效
	return tokenVersion >= currentVersion
}

