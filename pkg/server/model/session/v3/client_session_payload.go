package sessions

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ClientSessionPayload struct {
	ExpiresAt    int64  `json:"exp,omitempty"`           //expiration_time
	IssuedAt     int64  `json:"iat,omitempty"`           //issued_at_time
	Uuid         string `json:"uuid,omitempty"`          //token_uuid
	ClusterUuid  string `json:"cluster-uuid,omitempty"`  //cluster_uuid
	PollInterval int    `json:"poll-interval,omitempty"` //config_poll_interval
	Loglevel     string `json:"log-level,omitempty"`     //config_log_level
}

func (claims ClientSessionPayload) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().UTC().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !claims.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(claims.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("%s by %s", jwt.ErrTokenExpired, delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !claims.VerifyIssuedAt(now, false) {
		vErr.Inner = jwt.ErrTokenUsedBeforeIssued
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (claims ClientSessionPayload) VerifyExpiresAt(cmp int64, req bool) bool {
	if claims.ExpiresAt == 0 {
		return verifyExp(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(claims.ExpiresAt, 0)
	return verifyExp(&t, time.Unix(cmp, 0), req)
}

func (claims ClientSessionPayload) VerifyIssuedAt(cmp int64, req bool) bool {
	if claims.IssuedAt == 0 {
		return verifyIat(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(claims.IssuedAt, 0)
	return verifyIat(&t, time.Unix(cmp, 0), req)
}

func verifyExp(exp *time.Time, now time.Time, required bool) bool {
	if exp == nil {
		return !required
	}
	return now.Before(*exp)
}

func verifyIat(iat *time.Time, now time.Time, required bool) bool {
	if iat == nil {
		return !required
	}
	return now.After(*iat) || now.Equal(*iat)
}
