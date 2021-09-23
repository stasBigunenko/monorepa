package httphandler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/customErrors"
	mocks "github.com/stasBigunenko/monorepa/mocks/pkg/http/handler"
	"github.com/stasBigunenko/monorepa/model"
)

const (
	headerString string = "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzA1ODI1NzIsImlhdCI6MTYzMDU4MDc3MiwiSWQiOiI3Yzc2NTBmZS04NDNjLTQ3NmUtODEzMi1jZTc1NGUxNTMxNGMiLCJlbWFpbCI6IkJvYiJ9"
)

type MockTokenService struct {
}

func (s MockTokenService) ParseToken(tokenHeader string) (string, error) {
	return "", nil
}

type MockLoggingService struct {
}

func (s MockLoggingService) WriteLog(ctx context.Context, message string) {}

func TestHTTPHandler(t *testing.T) {
	type fields struct {
		AccountsService AccountGrpcService
		UsersService    UserGrpcService
	}

	type args struct {
		url     string
		method  string
		header  string
		body    []byte
		context string
	}

	type resp struct {
		code int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   resp
	}{ /*----------Testing Accounts---------*/
		{
			name: "POST /accounts OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockCreateAccount: func(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
						return uuid.New(), nil
					},
				},
			},
			args: args{
				url:    "/accounts",
				method: "POST",
				body:   []byte(`{"id":"32b56c48-1b96-11ec-adc6-23ffd7a72bbb"}`),
			},
			want: resp{code: http.StatusCreated},
		},
		{
			name: "POST /accounts !OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockCreateAccount: func(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
						return uuid.Nil, customErrors.AlreadyExists
					},
				},
			},
			args: args{
				url:    "/accounts",
				method: "POST",
				body:   []byte(`{"id":"32b56c48-1b96-11ec-adc6-23ffd7a72bbb"}`),
			},
			want: resp{code: http.StatusBadRequest},
		},
		{
			name: "GET /accounts/{id} OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetAccount: func(_ context.Context, _ uuid.UUID) (model.Account, error) {
						return model.Account{}, nil
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "GET /accounts/{id} !OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetAccount: func(_ context.Context, _ uuid.UUID) (model.Account, error) {
						return model.Account{}, customErrors.NotFound
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusNotFound},
		},
		{
			name: "GET /accounts OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetAllAccounts: func(_ context.Context) ([]model.Account, error) {
						return []model.Account{}, nil
					},
				},
			},
			args: args{
				url:    "/accounts",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "GET /accounts !OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetAllAccounts: func(_ context.Context) ([]model.Account, error) {
						return nil, errors.New("strange error")
					},
				},
			},
			args: args{
				url:    "/accounts",
				method: "GET",
			},
			want: resp{code: http.StatusInternalServerError},
		},
		{
			name: "PUT /accounts/{id} OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockUpdateAccount: func(_ context.Context, _ model.Account) error {
						return nil
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "PUT",
				body:   []byte(`{"balance":100}`),
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "PUT /accounts/{id} !OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockUpdateAccount: func(_ context.Context, _ model.Account) error {
						return customErrors.DeadlineExceeded
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "PUT",
				body:   []byte(`{"balance":100}`),
			},
			want: resp{code: http.StatusGatewayTimeout},
		},
		{
			name: "DELETE /accounts/{id} OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockDeleteAccount: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "DELETE",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "DELETE /accounts/{id} !OK",
			fields: fields{
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockDeleteAccount: func(_ context.Context, _ uuid.UUID) error {
						return customErrors.NotFound
					},
				},
			},
			args: args{
				url:    "/accounts/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "DELETE",
			},
			want: resp{code: http.StatusNotFound},
		},
		/*----------Testing Users---------*/
		{
			name: "POST /users OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockCreateUser: func(_ context.Context, _ string) (uuid.UUID, error) {
						return uuid.New(), nil
					},
				},
			},
			args: args{
				url:    "/users",
				method: "POST",
				body:   []byte(`{"id":"32b56c48-1b96-11ec-adc6-23ffd7a72bbb"}`),
			},
			want: resp{code: http.StatusCreated},
		},
		{
			name: "POST /users !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockCreateUser: func(_ context.Context, _ string) (uuid.UUID, error) {
						return uuid.Nil, customErrors.AlreadyExists
					},
				},
			},
			args: args{
				url:    "/users",
				method: "POST",
				body:   []byte(`{"id":"32b56c48-1b96-11ec-adc6-23ffd7a72bbb"}`),
			},
			want: resp{code: http.StatusBadRequest},
		},
		{
			name: "GET /users/{id} OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetUser: func(_ context.Context, _ uuid.UUID) (model.UserHTTP, error) {
						return model.UserHTTP{}, nil
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "GET /users/{id} !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetUser: func(_ context.Context, _ uuid.UUID) (model.UserHTTP, error) {
						return model.UserHTTP{}, customErrors.NotFound
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusNotFound},
		},
		{
			name: "GET /users OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetAllUsers: func(_ context.Context) ([]model.UserHTTP, error) {
						return []model.UserHTTP{}, nil
					},
				},
			},
			args: args{
				url:    "/users",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "GET /users !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetAllUsers: func(_ context.Context) ([]model.UserHTTP, error) {
						return nil, errors.New("strange error")
					},
				},
			},
			args: args{
				url:    "/users",
				method: "GET",
			},
			want: resp{code: http.StatusInternalServerError},
		},
		{
			name: "PUT /users/{id} OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockUpdateUser: func(_ context.Context, _ model.UserHTTP) error {
						return nil
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "PUT",
				body:   []byte(`{"name":"john"}`),
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "PUT /users/{id} !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockUpdateUser: func(_ context.Context, _ model.UserHTTP) error {
						return customErrors.DeadlineExceeded
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "PUT",
				body:   []byte(`{"name":"john"}`),
			},
			want: resp{code: http.StatusGatewayTimeout},
		},
		{
			name: "DELETE /users/{id} OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockDeleteUser: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "DELETE",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "DELETE /users/{id} !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockDeleteUser: func(_ context.Context, _ uuid.UUID) error {
						return customErrors.NotFound
					},
				},
			},
			args: args{
				url:    "/users/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "DELETE",
			},
			want: resp{code: http.StatusNotFound},
		}, /*----------Testing Aggregate---------*/
		{
			name: "GET /accounts_and_user/{id} OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetUser: func(_ context.Context, _ uuid.UUID) (model.UserHTTP, error) {
						return model.UserHTTP{}, nil
					},
				},
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetUserAccounts: func(_ context.Context, _ uuid.UUID) ([]model.Account, error) {
						return []model.Account{}, nil
					},
				},
			},
			args: args{
				url:    "/accounts_and_user/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "GET /accounts_and_user/{id} !OK",
			fields: fields{
				UsersService: &mocks.MockUsersGrpcServer{
					MockGetUser: func(_ context.Context, _ uuid.UUID) (model.UserHTTP, error) {
						return model.UserHTTP{}, customErrors.NotFound
					},
				},
				AccountsService: &mocks.MockAccountsGrpcServer{
					MockGetUserAccounts: func(_ context.Context, _ uuid.UUID) ([]model.Account, error) {
						return []model.Account{}, nil
					},
				},
			},
			args: args{
				url:    "/accounts_and_user/32b56c48-1b96-11ec-adc6-23ffd7a72bbb",
				method: "GET",
			},
			want: resp{code: http.StatusNotFound},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HTTPHandler{
				AccountsService: tt.fields.AccountsService,
				UsersService:    tt.fields.UsersService,
				TokenService:    MockTokenService{},
				LoggingService:  MockLoggingService{},
			}

			hs := httptest.NewServer(s.GetRouter())
			defer hs.Close()

			cl := hs.Client()
			req, _ := http.NewRequest(tt.args.method, hs.URL+tt.args.url, bytes.NewReader(tt.args.body))
			req.Header.Set(
				"Authorization", headerString,
			)

			r, err := cl.Do(req)

			if err != nil || r.StatusCode != tt.want.code {
				log.Print(r.StatusCode)
				if err != nil {
					t.Errorf("error: %s", err)
				} else {
					t.Errorf("%s %s = %v, want %v", tt.args.method, tt.args.url, r.StatusCode, tt.want.code)
				}
			}
		})
	}
}
