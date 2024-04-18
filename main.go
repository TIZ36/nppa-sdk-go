package main

import (
	"context"
	"net/http"

	"github.com/tiz36/nppa-sdk-go/domain"
	nppaapis "github.com/tiz36/nppa-sdk-go/nppa_apis"
)

func main() {
	// 2 steps to do nppa apis
	// 1. new nppa api
	apis := nppaapis.NewNppaApi(domain.NPPAEndpointConfig{}, nil)

	// 2. do request
	apis.PlayerBehaviorDataReport(
		context.Background(),
		[]domain.NPPAPlayerBehaviorReportCollection{},
		&http.Client{}, // http client should be re-used, may cause memory leak if not
	)
}
