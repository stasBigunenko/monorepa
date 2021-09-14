package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRSAPrivateKey(t *testing.T) {

	type request struct {
		certVersion string
		pathCerts   string
	}

	type want struct {
		isError bool
	}

	testCases := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "normal request",
			request: request{
				certVersion: "1",
				pathCerts:   "../../../pkg/storage/certificates",
			},
			want: want{isError: false},
		},
		{
			name: "wrong path",
			request: request{
				certVersion: "1",
				pathCerts:   "./",
			},
			want: want{
				isError: true,
			},
		},
		{
			name: "wrong cert version",
			request: request{
				certVersion: "0",
				pathCerts:   "../../../pkg/storage/certificates",
			},
			want: want{
				isError: true,
			},
		},
	}

	for _, tc := range testCases {
		privKey, err := readRSAPrivateKey(tc.request.certVersion, tc.request.pathCerts)

		if tc.want.isError {
			assert.Nil(t, privKey, tc.name)
			assert.Error(t, err, tc.name)
		} else {
			assert.NotNil(t, privKey, tc.name)
			assert.Nil(t, err, tc.name)
		}
	}
}
