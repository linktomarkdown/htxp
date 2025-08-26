package htxp

import (
	"crypto/md5"
	crand "crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"context"
	"net/http"
)

type RedisOptions redis.Options
type CacheClient struct {
	*redis.Client
}
type RegisterOptions struct {
	Name           string
	Type           string
	Path           string
	Command        string
	CommandStop    string
	Status         int
	RemoteHost     string
	RemotePort     int
	RemoteUser     string
	RemotePassword string
	RemoteKey      string
	ApiUrl         string
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// 业务标识映射
var businessCode = map[string]string{
	"alipay": "ALI",
	"wechat": "WX",
	"union":  "UN",
}

// 订单序列号（每秒递增）
var mutex sync.Mutex
var seqMap = make(map[string]int)

// GenerateOrderID 生成订单号（符合微信 32 位长度）
func GenerateOrderID(paymentType string) string {
	now := time.Now()

	// 1. 时间戳（14 位）
	timestamp := now.Format("20060102150405") // 精确到秒

	// 2. 业务标识（3 位）
	bizCode, exists := businessCode[paymentType]
	if !exists {
		bizCode = "UNK" // 默认未知业务
	}

	// 3. 序列号（4 位）- 保证同秒内的唯一性
	mutex.Lock()
	key := now.Format("20060102150405") // 以秒为单位
	seqMap[key]++
	seq := seqMap[key] % 10000 // 限制 4 位（0000 ~ 9999）
	mutex.Unlock()

	// 4. 随机数（6 位）- 防止并发冲突
	randomCode := fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))

	// 组合订单号
	return fmt.Sprintf("%s%s%04d%s", timestamp, bizCode, seq, randomCode)
}

// GenerateName 生成名称
func GenerateName(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

// Md5V 密码md5加密
func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRound 四舍五入
func GetRound(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func NewRedisConnect(options *RedisOptions) (*CacheClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &CacheClient{rdb}, nil
}

// StringToFloat64 字符串转浮点数
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// StringToInt 字符串转整数
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// InArray 判断元素是否在数组中
func InArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

// TryCatch 捕获异常
func TryCatch(f func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			handler(err)
		}
	}()
	f()
}

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[r.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[r.Intn(len(numBytes))]
		} else {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		}
	}
	return string(b)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// GenerateRandomNumber 生成随机数字
func GenerateRandomNumber(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = numBytes[r.Intn(len(numBytes))]
	}
	return string(b)
}

// GenerateRandomSpecial 生成随机特殊字符
func GenerateRandomSpecial(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = specialBytes[r.Intn(len(specialBytes))]
	}
	return string(b)
}

// GenerateRandomMixed 生成随机混合字符串
func GenerateRandomMixed(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	allBytes := letterBytes + numBytes + specialBytes
	for i := range b {
		b[i] = allBytes[r.Intn(len(allBytes))]
	}
	return string(b)
}

// AddPrefix 添加前缀
func AddPrefix(path string, prefix string) string {
	return prefix + path
}

// inArray 判断元素是否在数组中（私有方法）
func inArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// CopyDir 复制目录
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// ConvertUidToUint64 转换UID为uint64
func ConvertUidToUint64(uid string) uint64 {
	u, _ := strconv.ParseUint(uid, 10, 64)
	return u
}

// Contains 判断是否包含
func Contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

// GenerateKey 生成密钥
func GenerateKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func ConvertToMap(str string) map[string]string {
	var resultMap = make(map[string]string)
	values := strings.Split(str, "&")
	for _, value := range values {
		vs := strings.Split(value, "=")
		resultMap[vs[0]] = vs[1]
	}
	return resultMap
}

func GetUserIdAsUint64(data map[string]interface{}) (uint64, error) {
	// 检查 userId 是否存在
	userIdInterface, exists := data["userId"]
	if !exists {
		return 0, fmt.Errorf("userId not found in the message")
	}

	var userId uint64
	var err error

	// 根据不同类型进行处理
	switch v := userIdInterface.(type) {
	case float64:
		userId = uint64(v)
	case string:
		// 尝试将字符串转换为 uint64
		var num uint64
		_, err = fmt.Sscanf(v, "%d", &num)
		if err == nil {
			userId = num
		}
	default:
		err = fmt.Errorf("unsupported type for userId: %T", v)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to convert userId to uint64: %w", err)
	}

	return userId, nil
}

func GetUIDFromContext(r *http.Request) (uint64, error) {
	// 1. 获取用户 UID
	uid, ok := r.Context().Value("payload").(string)
	if !ok {
		return 0, fmt.Errorf("未获取到用户 UID")
	}
	// 2. 转换为 uint64
	var uidInt uint64
	_, err := fmt.Sscanf(uid, "%d", &uidInt)
	if err != nil {
		return 0, fmt.Errorf("failed to convert uid to uint64: %w", err)
	}
	return uidInt, nil
}

func GetUIDFromLogic(ctx context.Context) (uint64, error) {
	// 1. 获取用户 UID
	uid, ok := ctx.Value("payload").(string)
	if !ok {
		return 0, fmt.Errorf("未获取到用户 UID")
	}
	// 2. 转换为 uint64
	var uidInt uint64
	_, err := fmt.Sscanf(uid, "%d", &uidInt)
	if err != nil {
		return 0, fmt.Errorf("failed to convert uid to uint64: %w", err)
	}
	return uidInt, nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Paginate 分页工具
func Paginate(page, pageSize int64) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset = int((page - 1) * pageSize)
	limit = int(pageSize)
	return
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total, pageSize int64) int64 {
	if total == 0 {
		return 0
	}
	return int64(math.Ceil(float64(total) / float64(pageSize)))
}

// StringPtr 字符串指针
func StringPtr(s string) *string {
	return &s
}

// Int64Ptr 整数指针
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr 浮点数指针
func Float64Ptr(f float64) *float64 {
	return &f
}

// NullFloat64Ptr 创建 sql.NullFloat64
func NullFloat64Ptr(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

// NullInt64Ptr 创建 sql.NullInt64
func NullInt64Ptr(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// NullStringPtr 创建 sql.NullString
func NullStringPtr(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NullTimePtr 创建 sql.NullTime
func NullTimePtr(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
