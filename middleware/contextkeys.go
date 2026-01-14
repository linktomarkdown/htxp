package middleware

// contextKey 用于 context 的 key，避免使用 string 类型导致冲突
type contextKey string

const (
	// ContextKeyUserID 用户ID
	ContextKeyUserID contextKey = "userId"
	// ContextKeyIsAdmin 是否是管理员
	ContextKeyIsAdmin contextKey = "isAdmin"
	// ContextKeyUserInfo 用户完整信息（可选）
	ContextKeyUserInfo contextKey = "userInfo"
	// ContextKeyIsTutorialVip 是否是教程会员
	ContextKeyIsTutorialVip contextKey = "isTutorialVip"
	// ContextKeyPurchasedLessonIds 已购买的课程ID列表
	ContextKeyPurchasedLessonIds contextKey = "purchasedLessonIds"
	// ContextKeyIsMonthlyVip 是否是月度会员
	ContextKeyIsMonthlyVip contextKey = "isMonthlyVip"
	// ContextKeyIsYearlyVip 是否是年度会员
	ContextKeyIsYearlyVip contextKey = "isYearlyVip"
)
