package domain

import "context"

// NPPA 领域接口
// 主要定义了NPPA领域应该提供的功能方法
type NPPA interface {
	// RealNameAuth 实名认证
	RealNameAuth(
		ctx context.Context,
		internalUid string,
		realName string,
		idNum string,
		clientOptions ...NppaClientOption,
	) (rns RealNameAuthStatus, personId string, err error)

	// RealNameAuthQuery 实名认证查询
	RealNameAuthQuery(
		ctx context.Context,
		internalUid string,
		clientOptions ...NppaClientOption,
	) (rns RealNameAuthStatus, personId string, err error)

	// PlayerBehaviorDataReport 游戏玩家行为数据上报
	PlayerBehaviorDataReport(
		ctx context.Context,
		collections []NPPAPlayerBehaviorReportCollection,
		clientOptions ...NppaClientOption,
	) (behaviorReportResults []NPPAPlayerBehaviorDataReportResult, err error)
}
