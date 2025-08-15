package v1

type DailyResp struct{}

// CalendarReq 签到日历请求参数
type CalendarReq struct {
	YearMonth string `form:"yearMonth" binding:"required"` // ?yearMonth=2025-01
}

// CalendarResp 签到日历响应参数
type CalendarResp struct {
	Year   int        `json:"year"`
	Month  int        `json:"month"`
	Detail DetailInfo `json:"detail"`
}

type DetailInfo struct {
	CheckedInDays      []int `json:"checkedInDays"`      // 签到的日期序号
	RetroCheckedInDays []int `json:"retroCheckedInDays"` // 补签的日期序号
	IsCheckedInToday   bool  `json:"isCheckedInToday"`   // 今天是否签到
	RemainRetroTimes   int   `json:"remainRetroTimes"`   // 剩余补签次数
	ConsecutiveDays    int   `json:"consecutiveDays"`    // 连续签到天数
}

// RetroReq 补签请求参数
type RetroReq struct {
	Date string `json:"date" binding:"required"`
}

// RetroResp 补签响应参数
type RetroResp struct{}
