package oauth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/luizmoitinho/bookstore_oauth/errors"
	"github.com/luizmoitinho/bookstore_oauth/rest"
	"gopkg.in/resty.v1"
)

const (
	HEADER_X_PUBLIC    = "X-Public"
	HEADER_X_CLIENT_ID = "X-Client-Id"
	HEADER_X_CALLER_ID = "X-Caller-Id"

	PARAM_ACCESS_TOKEN = "access_token"

	ACCESS_TOKEN_BASE_URL = "http://localhost:8080"
	ACCESS_TOKEN_URI      = "/oauth/access_token/%s"
)

func IsPublic(in *http.Request) bool {
	if in == nil {
		return false
	}

	return in.Header.Get(HEADER_X_PUBLIC) == "true"
}

func GetCallerId(in *http.Request) int64 {
	if in == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(in.Header.Get(HEADER_X_CALLER_ID), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func GetClientId(in *http.Request) int64 {
	if in == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(in.Header.Get(HEADER_X_CLIENT_ID), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func Authenticate(in *http.Request) *errors.Rest {
	if in == nil {
		return nil
	}

	cleanRequest(in)

	accessToken := strings.TrimSpace(in.URL.Query().Get(PARAM_ACCESS_TOKEN))
	if accessToken == "" {
		return nil
	}

	client := rest.NewClient(resty.New(), ACCESS_TOKEN_BASE_URL, ACCESS_TOKEN_URI)
	at, err := client.GetAccessToken(accessToken)
	if err != nil {
		if err.Status == http.StatusNotFound {
			return nil
		}
		return err
	}

	in.Header.Add(HEADER_X_CLIENT_ID, fmt.Sprintf("%v", at.ClientId))
	in.Header.Add(HEADER_X_CALLER_ID, fmt.Sprintf("%v", at.UserId))

	return nil
}

func cleanRequest(in *http.Request) {
	if in == nil {
		return
	}
	in.Header.Del(HEADER_X_CLIENT_ID)
	in.Header.Del(HEADER_X_CALLER_ID)
}
