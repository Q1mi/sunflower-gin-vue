package v1

type SummaryResp struct {
	Total int64 `json:"total"`
}

// RecordsReq 积分记录请求结构体
type RecordsReq struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// RecordsResp 积分记录响应结构体
type RecordsResp struct {
	Total   int64         `json:"total"`
	HasMore bool          `json:"hasMore"` // 是否还有更多数据
	List    []*RecordInfo `json:"list"`
}

type RecordInfo struct {
	PointsChange    int64  `json:"pointsChange"`
	TransactionType int32  `json:"transactionType"`
	Description     string `json:"description"`
	TransactionTime string `json:"transactionTime"`
}
