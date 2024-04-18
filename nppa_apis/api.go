package nppaapis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tiz36/nppa-sdk-go/domain"
)

type NppaApi struct {
	Host          string
	Api           string
	NppaAppId     string
	NppaSecretKey string
	BizId         string

	Logger domain.Logger
}

func NewNppaApi(nppaConfig domain.NPPAEndpointConfig, logger domain.Logger) *NppaApi {
	return &NppaApi{
		Host:          nppaConfig.Host,
		Api:           nppaConfig.Api,
		NppaAppId:     nppaConfig.NppaAppId,
		NppaSecretKey: nppaConfig.NppaSecretKey,
		BizId:         nppaConfig.BizId,
		Logger:        logger,
	}
}

// RealNameAuth 实名认证
func (nppaApi *NppaApi) RealNameAuth(
	ctx context.Context,
	internalUid string,
	realName string,
	idNum string,
	httpClient *http.Client,
	clientOptions ...domain.NppaClientOption,
) (rns domain.RealNameAuthStatus, personId string, err error) {
	var realNameAuthResp domain.NPPARealNameAuthResp
	var body []byte
	var nppaHttpRawResp []byte
	var nppaClient *domain.NppaHttpClient

	// 1. 初始化NPPA客户端
	if nppaClient, err = domain.NewNppaClient(
		nppaApi.Host,
		nppaApi.NppaAppId,
		nppaApi.NppaSecretKey,
		nppaApi.BizId,
		httpClient,
		clientOptions...,
	); err != nil {
		return domain.RealNameAuthStatusFail, "", err
	}

	// 请求参数
	if body, err = json.Marshal(domain.NPPARealNameAuthRequest{
		AI:    internalUid,
		Name:  realName,
		IdNum: idNum,
	}); err != nil {
		return domain.RealNameAuthStatusFail, "", err
	}

	nppaApi.Logger.Infof("RealNameAuth, Req =-===> %s", string(body))

	// 发送请求
	if nppaHttpRawResp, err = nppaClient.Request(ctx, http.MethodPost, nppaApi.Api, nil, body); err != nil {
		return domain.RealNameAuthStatusFail, "", err
	}

	nppaApi.Logger.Infof("RealNameAuth, Resp =-===> %s", string(nppaHttpRawResp))

	// 解析响应
	err = json.Unmarshal(nppaHttpRawResp, &realNameAuthResp)
	if err != nil {
		return domain.RealNameAuthStatusFail, "", err
	}

	if realNameAuthResp.ErrCode != 0 {
		return domain.RealNameAuthStatusFail, "", fmt.Errorf("errcode: %d, errmsg: %s", realNameAuthResp.ErrCode, realNameAuthResp.ErrMsg)
	}

	// 返回结果
	return realNameAuthResp.Data.Result.Status, realNameAuthResp.Data.Result.Pi, nil
}

// RealNameAuthQuery 实名认证查询
func (nppaApi *NppaApi) RealNameAuthQuery(
	ctx context.Context,
	internalUid string,
	httpClient *http.Client,
	clientOptions ...domain.NppaClientOption,
) (rns domain.RealNameAuthStatus, personId string, err error) {
	var realNameAuthQueryResp domain.NPPARealNameAuthQueryResp
	var values = url.Values{}
	var nppaHttpRawResp []byte
	var nppaClient *domain.NppaHttpClient

	// 1. 初始化NPPA客户端
	if nppaClient, err = domain.NewNppaClient(nppaApi.Host, nppaApi.NppaAppId, nppaApi.NppaSecretKey, nppaApi.BizId, httpClient, clientOptions...); err != nil {
		return domain.RealNameAuthStatusFail, "", err
	}

	// 请求参数
	values.Add("ai", internalUid)

	nppaApi.Logger.Infof("RealNameAuthQuery, req =-===> %s", values.Encode())

	// 发送请求
	if nppaHttpRawResp, err = nppaClient.Request(ctx, http.MethodGet, nppaApi.Api, values, nil); err != nil {
		return domain.RealNameAuthStatusProcessing, "", err
	}

	nppaApi.Logger.Infof("RealNameAuthQuery, resp =-===> %s", string(nppaHttpRawResp))

	// 解析响应
	err = json.Unmarshal(nppaHttpRawResp, &realNameAuthQueryResp)
	if err != nil {
		return domain.RealNameAuthStatusProcessing, "", err
	}

	if realNameAuthQueryResp.ErrCode != 0 {
		return domain.RealNameAuthStatusProcessing, "", fmt.Errorf("errcode: %d, errmsg: %s", realNameAuthQueryResp.ErrCode, realNameAuthQueryResp.ErrMsg)
	}

	// 返回结果
	return realNameAuthQueryResp.Data.Result.Status, realNameAuthQueryResp.Data.Result.Pi, nil
}

// PlayerBehaviorDataReport 游戏玩家行为数据上报
func (nppaApi *NppaApi) PlayerBehaviorDataReport(
	ctx context.Context,
	collections []domain.NPPAPlayerBehaviorReportCollection,
	httpClient *http.Client,
	clientOptions ...domain.NppaClientOption,
) (behaviorReportResults []domain.NPPAPlayerBehaviorDataReportResult, err error) {

	var resp domain.NPPAPlayerBehaviorDataReportResp
	var body []byte
	var nppaHttpRawResp []byte
	var nppaClient *domain.NppaHttpClient

	// 1. 初始化NPPA客户端
	if nppaClient, err = domain.NewNppaClient(nppaApi.Host, nppaApi.NppaAppId, nppaApi.NppaSecretKey, nppaApi.BizId, httpClient, clientOptions...); err != nil {
		return nil, err
	}

	// 1. 请求参数
	if body, err = json.Marshal(domain.NPPAPlayerBehaviorDataReportRequest{
		Collections: collections,
	}); err != nil {
		return nil, err
	}

	nppaApi.Logger.Infof("UserBehaviorReport, req =-===> %s", string(body))

	// 2. 发送请求
	if nppaHttpRawResp, err = nppaClient.Request(ctx, http.MethodPost, nppaApi.Api, nil, body); err != nil {
		return nil, err
	}

	nppaApi.Logger.Infof("UserBehaviorReport, req =-===> %s", string(nppaHttpRawResp))

	// 3. 解析响应
	if err = json.Unmarshal(nppaHttpRawResp, &resp); err != nil {
		return nil, err
	}

	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("errcode: %d, errmsg: %s", resp.ErrCode, resp.ErrMsg)
	}

	return resp.Data.Results, nil
}
