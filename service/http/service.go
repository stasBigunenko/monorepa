package httpservice

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	jwtservice "github.com/stasBigunenko/monorepa/service/auth/jwt"
)

type HTTPService struct {
	JwtServiceAddr string
}

type tokenResp struct {
	PbKey []byte `json:"publicKey"`
}

func (s HTTPService) ParseToken(tokenHeader string) error {
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return fmt.Errorf("malformed auth token, could not split two parts")
	}

	if splitted[0] != "bearer" {
		return fmt.Errorf("malformed auth token, the first part is not bearer")

	}

	tokenPart := splitted[1]

	token, err := jwt.ParseWithClaims(tokenPart, &jwtservice.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwtservice.UserClaims)
		if !ok {
			return nil, fmt.Errorf("wrong format of claims")
		}

		expTime := time.Unix(claims.ExpiresAt, 0)
		if !expTime.After(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}

		resp, err := http.Get(path.Join(s.JwtServiceAddr, claims.KeyVersion))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to jwt server: %w", err)
		}
		defer resp.Body.Close()

		var publickeyJSON tokenResp
		if err = json.NewDecoder(resp.Body).Decode(&publickeyJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
		}

		block, _ := pem.Decode(publickeyJSON.PbKey)
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

		publicKeyForToken, ok := pubInterface.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("failed to convert public key: %w", err)
		}

		return publicKeyForToken, nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("token is not valid: %w", err)
	}

	return nil
}
