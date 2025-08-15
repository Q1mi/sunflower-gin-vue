package model

type MonthDetailOutput struct {
	CheckedInDays      []int `json:"checkedInDays"`      // 签到的日期序号
	RetroCheckedInDays []int `json:"retroCheckedInDays"` // 补签的日期序号
	IsCheckedInToday   bool  `json:"isCheckedInToday"`   // 今天是否签到
	RemainRetroTimes   int   `json:"remainRetroTimes"`   // 剩余补签次数
	ConsecutiveDays    int   `json:"consecutiveDays"`    // 连续签到天数
}
