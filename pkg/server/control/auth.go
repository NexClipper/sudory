package control

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	auth "github.com/NexClipper/sudory/pkg/server/model/auth/v2"
	tenantv3 "github.com/NexClipper/sudory/pkg/server/model/tenant/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description auth
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/auth
// @Router      /server/auth [post]
// @Param       object       body   v2.HttpReq_ServiceAccessToken true  "v2.HttpReq_ServiceAccessToken"
// @Success     200 {object} v2.HttpRsp_AccessTokenResponse
func (ctl ControlVanilla) Auth(ctx echo.Context) (err error) {
	body := new(auth.HttpReq_ServiceAccessToken)
	err = func() (err error) {
		err = ctx.Bind(body)
		// err = echoutil.Bind(ctx, body)
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// trim
	pattern := strings.TrimSpace(body.Tenant)
	hash := SHA1(pattern)
	var id int64

	// get tenant
	tenant := tenantv3.Tenant{
		Hash: hash,
	}
	tenant_cond := stmt.Equal("hash", tenant.Hash)
	err = stmtex.Select(tenant.TableName(), []string{"id"}, tenant_cond, nil, nil).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, i int) error {
			return scan.Scan(&id)
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get tenant")
	}

	time_now := time.Now()

	if id == 0 {
		// save tenant
		new_tenant := tenantv3.NewTenant(hash, pattern, time_now)
		updatecolumns := []string{"hash"}
		save_tenant := stmtex.InsertOrUpdate(new_tenant.TableName(), new_tenant.ColumnNames(), updatecolumns, new_tenant.Values())
		_, lastid, err := save_tenant.ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "failed to save new tenant")
		}
		id = lastid // get tenant id
	}

	claims := auth.ServiceAccessTokenClaims{
		ID:        id,
		Hash:      hash,
		Tenant:    pattern,
		IssuedAt:  time_now.Unix(),
		ExpiresAt: globvar.ServiceSession.ExpirationTime(time_now).Unix(),
	}

	token_str, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(globvar.ServiceSession.SignatureSecret()))

	if err != nil {
		return errors.Wrapf(err, "failed to create new service access token")
	}

	return ctx.JSON(http.StatusOK, auth.HttpRsp_AccessTokenResponse{
		TokenType:   "Bearer",
		AccessToken: token_str,
		ExpiresIn:   int(claims.ExpiresAt - claims.IssuedAt),
	})
}

func SHA1(s string) string {
	type bytes = []byte

	h := sha1.Sum(bytes(s))
	x := hex.EncodeToString(h[0:])
	return x
}
