package model

// AddPointInput 添加积分输入参数
type AddPointInput struct {
	UserID      int64
	PointAmount int64
	Type        int32
	Desc        string
}

type SummaryOutput struct {
	TotalPoint int64
}

type RecordsInput struct {
	UserID int64
	Offset int
	Limit  int
}

type RecordsOutput struct {
	Total   int64
	HasMore bool
	List    []*RecordInfo
}

type RecordInfo struct {
	PointsChange    int64
	TransactionType int32
	Description     string
	TransactionTime string
}
