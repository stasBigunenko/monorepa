package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"monorepa/model"
	"monorepa/service/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserTokenGen(t *testing.T) {
	// test info
	type request struct {
		endpoint string
		method   string
		body     func() []byte
	}

	type wantResp struct {
		code        int
		headerToken bool
		body        func() []byte
	}

	testCases := []struct {
		name            string
		service         *auth.MockService // service to mock
		serviceFuncResp func(*auth.MockService, []byte)
		request         request
		want            wantResp
	}{
		{
			name:    "Normal request",
			service: &auth.MockService{},
			serviceFuncResp: func(mc *auth.MockService, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("qwerty.qwerty.qwerty", nil)
			},
			request: request{
				endpoint: "/login",
				method:   "GET",
				body: func() []byte {
					user := model.User{
						Name:     "Bob",
						Password: "12345",
					}
					res, _ := json.Marshal(user)
					return res
				},
			},
			want: wantResp{
				code: 201,
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},
				headerToken: true,
			},
		},
		{
			name:    "Wirhout user data",
			service: &auth.MockService{},
			serviceFuncResp: func(mc *auth.MockService, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("", errors.New("wrong password"))
			},
			request: request{
				endpoint: "/login",
				method:   "GET",
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},
			},
			want: wantResp{
				code: 400,
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},
				headerToken: false,
			},
		},
		{
			name:    "Wirhout user password",
			service: &auth.MockService{},
			serviceFuncResp: func(mc *auth.MockService, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("", errors.New("wrong password"))
			},
			request: request{
				endpoint: "/login",
				method:   "GET",
				body: func() []byte {
					user := model.User{
						Name: "Bob",
					}
					res, _ := json.Marshal(user)
					return res
				},
			},
			want: wantResp{
				code: 400,
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},

				headerToken: false,
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := HandlerItemsServ{
			router:   mux.NewRouter(),
			ctx:      context.Background(),
			services: tc.service,
		}

		tc.serviceFuncResp(tc.service, tc.request.body()) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)

		token := res.Header.Get("token")
		if tc.want.headerToken {
			assert.NotEmpty(t, token, tc.name)
		} else {
			assert.Empty(t, token, tc.name)
		}
	}
}

func TestGetCertificateKey(t *testing.T) {

	type tokenResp struct {
		PbKey []byte `json:"publicKey"`
	}

	// test info
	type request struct {
		endpoint string
		method   string
	}

	type wantResp struct {
		code   int
		isBody bool
	}

	testCases := []struct {
		name            string
		service         *auth.MockService // service to mock
		serviceFuncResp func(*auth.MockService)
		request         request
		want            wantResp
	}{
		{
			name:    "Normal request",
			service: &auth.MockService{},
			serviceFuncResp: func(mc *auth.MockService) {
				mc.On("GetCert",
					mock.Anything,
				).Return([]byte("qwerty.qwerty.qwerty"), nil)
			},
			request: request{
				endpoint: "/get-cert/1",
				method:   "GET",
			},
			want: wantResp{
				code:   200,
				isBody: true,
			},
		},
		{
			name:    "Internal error",
			service: &auth.MockService{},
			serviceFuncResp: func(mc *auth.MockService) {
				mc.On("GetCert",
					mock.Anything,
				).Return([]byte{}, errors.New("some err"))
			},
			request: request{
				endpoint: "/get-cert/foo=bar&baz=bar",
				method:   "GET",
			},
			want: wantResp{
				code:   500,
				isBody: false,
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := HandlerItemsServ{
			router:   mux.NewRouter(),
			ctx:      context.Background(),
			services: tc.service,
		}

		tc.serviceFuncResp(tc.service) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, nil)

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)

		if tc.want.isBody {
			var respData tokenResp
			err = json.NewDecoder(res.Body).Decode(&respData)
			assert.Nil(t, err, tc.name)

			assert.NotEmpty(t, respData.PbKey, tc.name)
		} else {
			assert.Empty(t, res.Body, tc.name)
		}
	}
}
