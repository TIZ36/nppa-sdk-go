package domain

import (
	"net/url"
)

type RealNameAuthStatus int

type NPPAPlayerBehaviorType int

type NPPAPlayerCertType int

const (
	RealNameAuthStatusSuccess    RealNameAuthStatus = 0
	RealNameAuthStatusProcessing RealNameAuthStatus = 1
	RealNameAuthStatusFail       RealNameAuthStatus = 2

	NPPAPlayerBehaviorTypeLogin    NPPAPlayerBehaviorType = 1
	NPPAPlayerBehaviorTypeRegister NPPAPlayerBehaviorType = 0

	CertTypeAuthorized NPPAPlayerCertType = 0
	CertTypeVisitor    NPPAPlayerCertType = 2
)

type NPPATestEndpointConfig struct {
	RealNameAuth   NPPAEndpointConfig `json:"real_name_auth"`
	RealNameQuery  NPPAEndpointConfig `json:"real_name_query"`
	PlayerBehavior NPPAEndpointConfig `json:"player_behavior"`
}

type IdentityCsvData struct {
	AppId  string `json:"app_id"`
	AppUid string `json:"app_uid"`
	IdNum  string `json:"code"`
	Name   string `json:"name"`
	Pi     string `json:"pi"`
	Status string `json:"status"`
}

type NPPAUserInfo struct {
	AI    string `json:"ai"`
	Name  string `json:"name"`
	IdNum string `json:"idNum"`
}

type NPPAEndpointConfig struct {
	Host          string `json:"host"`
	Api           string `json:"api"`
	NppaAppId     string `json:"nppa_app_id"`
	NppaSecretKey string `json:"nppa_secret_key"`
	BizId         string `json:"biz_id"`

	// 这个配置项是用来跳过https证书校验的，因为中宣部测试的环境需要https调用，但是没有提供证书
	IsHttpsCertSkipVerify bool `json:"is_https_cert_skip_verify"`
}

// NPPACommonResp NPPA 通用响应
type NPPACommonResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// NPPACommonHeader NPPA 通用请求头
type NPPACommonHeader struct {
	ContentType string `json:"Content-Type"`
	AppId       string `json:"appId"`
	BizId       string `json:"bizId"`
	Timestamps  string `json:"timestamps"`
	Sign        string `json:"sign"`
}

type NPPAEncryptedReqBody struct {
	Data string `json:"data"`
}

// NPPARealNameAuthRequest
// 1 实名认证 请求、响应结构体
type NPPARealNameAuthRequest struct {
	AI    string `json:"ai"`
	Name  string `json:"name"`
	IdNum string `json:"idNum"`
}
type NPPARealNameAuthResult struct {
	Result struct {
		Status RealNameAuthStatus `json:"status"`
		Pi     string             `json:"pi"`
	} `json:"result"`
}
type NPPARealNameAuthResp struct {
	NPPACommonResp
	Data NPPARealNameAuthResult `json:"data"`
}

// NPPARealNameAuthQueryRequest
// 2 实名认证查询 请求、响应结构体
type NPPARealNameAuthQueryRequest = url.Values
type NPPARealNameAuthQueryResp struct {
	NPPACommonResp
	Data NPPARealNameAuthResult `json:"data"`
}

// NPPAPlayerBehaviorDataReportRequest
// 3 游戏玩家行为数据上报 请求、响应结构体
type NPPAPlayerBehaviorDataReportRequest struct {
	Collections []NPPAPlayerBehaviorReportCollection `json:"collections"`
}
type NPPAPlayerBehaviorDataReportResult struct {
	NPPACommonResp
	No int `json:"no"`
}
type NPPAPlayerBehaviorDataReportResults struct {
	Results []NPPAPlayerBehaviorDataReportResult `json:"results"`
}
type NPPAPlayerBehaviorDataReportResp struct {
	NPPACommonResp
	Data NPPAPlayerBehaviorDataReportResults `json:"data"`
}
type NPPAPlayerBehaviorReportCollection struct {
	No int    `json:"no"`
	Si string `json:"si"`
	Bt int    `json:"bt"`
	Ot int64  `json:"ot"`
	Ct int    `json:"ct"`
	Di string `json:"di"`
	Pi string `json:"pi"`
}
