package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authMock "github.com/stasBigunenko/monorepa/mocks/service/auth"
	"github.com/stasBigunenko/monorepa/model"

	"github.com/gorilla/mux"
	er "github.com/stasBigunenko/monorepa/customErrors"
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
		service         *authMock.Service // service to mock
		serviceFuncResp func(*authMock.Service, []byte)
		request         request
		want            wantResp
	}{
		{
			name:    "Normal request",
			service: &authMock.Service{},
			serviceFuncResp: func(mc *authMock.Service, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("qwerty.qwerty.qwerty", nil)
			},
			request: request{
				endpoint: "/login",
				method:   "POST",
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
			service: &authMock.Service{},
			serviceFuncResp: func(mc *authMock.Service, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("", er.WrongPassword)
			},
			request: request{
				endpoint: "/login",
				method:   "POST",
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
			name:    "Without user password",
			service: &authMock.Service{},
			serviceFuncResp: func(mc *authMock.Service, items []byte) {
				mc.On("Login",
					mock.Anything,
				).Return("", er.WrongPassword)
			},
			request: request{
				endpoint: "/login",
				method:   "POST",
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
		service         *authMock.Service // service to mock
		serviceFuncResp func(*authMock.Service)
		request         request
		want            wantResp
	}{
		{
			name:    "Normal request",
			service: &authMock.Service{},
			serviceFuncResp: func(mc *authMock.Service) {
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
			service: &authMock.Service{},
			serviceFuncResp: func(mc *authMock.Service) {
				mc.On("GetCert",
					mock.Anything,
				).Return([]byte{}, er.WrongPassword)
			},
			request: request{
				endpoint: "/get-cert/1",
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
