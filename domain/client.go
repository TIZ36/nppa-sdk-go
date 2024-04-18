package domain

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type NppaHttpClient struct {
	AppId        string
	BizId        string
	SecretKey    string
	SecretKeyHex []byte

	HttpClient *http.Client

	NppaHost string
}

type NppaRequestBodyData struct {
	Data string `json:"data"`
}

func (nppaReqBodyData *NppaRequestBodyData) String() string {
	if nppaReqBodyData == nil {
		return ""
	}

	bin, _ := json.Marshal(nppaReqBodyData)
	return string(bin)
}

func (nppaReqBodyData *NppaRequestBodyData) Byte() []byte {
	if nppaReqBodyData == nil {
		return nil
	}
	bin, _ := json.Marshal(nppaReqBodyData)
	return bin
}

type NppaClientOption func(*NppaHttpClient) (*NppaHttpClient, error)

func NewNppaClient(nppaHost, appId, secretKey string, bizId string, httpClient *http.Client, opts ...NppaClientOption) (*NppaHttpClient, error) {
	secretKeyHex, _ := hex.DecodeString(secretKey)

	defaultNppaClient := &NppaHttpClient{
		AppId:        appId,
		BizId:        bizId,
		SecretKey:    secretKey,
		SecretKeyHex: secretKeyHex,
		NppaHost:     nppaHost,
		HttpClient:   httpClient,
	}

	// apply options, 允许用户自定义配置
	var nppaClient = defaultNppaClient
	var err error
	for _, opt := range opts {
		if nppaClient, err = opt(nppaClient); err != nil {
			return nil, err
		}
	}

	return nppaClient, nil
}

func (nppaClient *NppaHttpClient) Request(
	ctx context.Context,
	httpMethod, api string,
	values url.Values,
	body []byte) ([]byte, error) {

	var requestBodyData *NppaRequestBodyData
	var signature string
	var err error

	// 1. 拼接请求URL
	baseUrl := nppaClient.NppaHost + api

	// 1.1 如果url.Values不为空，拼接请求URL参数
	if values != nil && len(values) > 0 {
		baseUrl += "?" + values.Encode()
	} else {
		values = make(url.Values)
	}

	// 2. 如果body不为空，加密请求体
	if body != nil && len(body) > 0 {
		// 2.1 加密请求体
		if requestBodyData, err = nppaClient.EncryptReBody(body); err != nil {
			return nil, err
		}
	} else {
		body = []byte{}
	}

	// 3. 创建请求
	req, err := http.NewRequest(httpMethod, baseUrl, bytes.NewBuffer(requestBodyData.Byte()))
	if err != nil {
		return nil, err
	}

	// 4. 设置NPPA特殊请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("appId", nppaClient.AppId)
	req.Header.Set("bizId", nppaClient.BizId)
	req.Header.Set("timestamps", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// 5. NPPA特殊请求头 签名
	// 5.1 拼接header到 url.Values, 字典序后字符串 SortedUrlValuesStr
	// 5.2 secretKeyStr+{SortedUrlValuesStr}+requestBodyDataStr 为待签名字符串
	if signature, err = nppaClient.DoSignature(values, requestBodyData.String()); err != nil {
		return nil, err
	}
	req.Header.Set("sign", signature)

	// 6. 发送请求
	resp, err := nppaClient.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (nppaClient *NppaHttpClient) DoSignature(values url.Values, requestBodyDataStr string) (string, error) {
	// 1. 拼接header到values
	values.Add("appId", nppaClient.AppId)
	values.Add("bizId", nppaClient.BizId)
	values.Add("timestamps", fmt.Sprintf("%d", time.Now().UnixMilli()))

	var paramList []string
	for k := range values {
		paramList = append(paramList, k+values.Get(k))
	}

	sort.Strings(paramList)
	var data = nppaClient.SecretKey + strings.Join(paramList, "") + requestBodyDataStr

	// 2. 计算签名 sha256
	sign := sha256.Sum256([]byte(data))
	hashedSign := hex.EncodeToString(sign[:])
	return hashedSign, nil
}

func (nppaClient *NppaHttpClient) EncryptReBody(reqBody []byte) (*NppaRequestBodyData, error) {
	aesGcm128KeyBin, _ := hex.DecodeString(nppaClient.SecretKey)
	block, err := aes.NewCipher(aesGcm128KeyBin)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encryptedBin := aesgcm.Seal(nonce, nonce, reqBody, nil)
	return &NppaRequestBodyData{
		Data: base64.StdEncoding.EncodeToString(encryptedBin),
	}, nil
}
