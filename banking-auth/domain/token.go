package domain

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Claims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}

func (c Claims) IsUserRole() bool {
	return c.Role == "user"
}

func BuildClaimsFromJwtMapClaims(mapClaims jwt.MapClaims) (*Claims, error) {
	data, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (c Claims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	//for admin role
	return true
}

func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}
	return true
}
