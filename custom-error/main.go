package main

import (
	"net/http"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/errorx"
)

var (
	// DOMAIN CUSTOM ERROR
	ERROR_INVALID_DATA = errorx.NewError(http.StatusBadRequest, "201", "Invalid data")
	ERROR_INFRA_ERROR  = errorx.NewError(http.StatusBadRequest, "202", "Infra Internal Error")
)
