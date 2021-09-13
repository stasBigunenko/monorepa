package jwt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRSAPrivateKey(t *testing.T) {

	type request struct {
		certVersion string
	}

	type want struct {
		isError bool
	}

	type env struct {
		set   func()
		unset func()
	}

	testCases := []struct {
		name    string
		request request
		want    want
		env     env
	}{
		{
			name: "normal request",
			request: request{
				certVersion: "1",
			},
			want: want{
				isError: false,
			},
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
			},
		},
		{
			name: "wrong path",
			request: request{
				certVersion: "1",
			},
			want: want{
				isError: true,
			},
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "./")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
			},
		},
		{
			name: "wrong cert version",
			request: request{
				certVersion: "0",
			},
			want: want{
				isError: true,
			},
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.env.set()
		defer tc.env.unset()

		privKey, err := readRSAPrivateKey(tc.request.certVersion)

		if tc.want.isError {
			assert.Nil(t, privKey, tc.name)
			assert.Error(t, err, tc.name)
		} else {
			assert.NotNil(t, privKey, tc.name)
			assert.Nil(t, err, tc.name)
		}
	}
}
