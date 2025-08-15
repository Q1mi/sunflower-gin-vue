package points

import (
	"context"
	"errors"
	"sunflower-gin/internal/dao/query"
	"sunflower-gin/internal/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Summary 查询用户积分信息
func Summary(ctx context.Context, userID int64) (*model.SummaryOutput, error) {
	output := &model.SummaryOutput{}
	// 1. 从数据库中查询用户积分信息
	upInst, err := query.UserPoint.WithContext(ctx).
		Where(query.UserPoint.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果是新用户还没有积分信息，则直接返回
			return output, nil
		}
		zap.L().Error("query user point error", zap.Error(err))
		return nil, err
	}
	// 2. 将数据封装到结构体中返回
	output.TotalPoint = upInst.Points
	return output, nil
}

// Records 查询用户积分记录
func Records(ctx context.Context, input *model.RecordsInput) (*model.RecordsOutput, error) {
	// 1. 从数据库中分页查询用户积分记录
	var records []*model.UserPointsTransaction
	total, err := query.UserPointsTransaction.WithContext(ctx).
		Where(query.UserPointsTransaction.UserID.Eq(input.UserID)).
		Order(query.UserPointsTransaction.CreatedAt.Desc()).
		ScanByPage(&records, input.Offset, input.Limit)
	if err != nil {
		zap.L().Error("query user points transaction error", zap.Error(err))
		return nil, err
	}
	// 2. 格式化数据
	list := make([]*model.RecordInfo, 0, len(records))
	for _, v := range records {
		list = append(list, &model.RecordInfo{
			PointsChange:    v.PointsChange,
			TransactionType: v.TransactionType,
			Description:     v.Description,
			TransactionTime: v.CreatedAt.Format(time.DateTime),
		})
	}
	hasMore := len(records) == input.Limit &&
		int(total) > input.Offset+input.Limit
	// 3. 返回数据
	return &model.RecordsOutput{
		Total:   total,
		HasMore: hasMore,
		List:    list,
	}, nil
}
