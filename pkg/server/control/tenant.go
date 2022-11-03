package control

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/model/auths/v2"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description tenant
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/tenant
// @Router      /server/tenant [post]
// @Param       object       body   auths.HttpReq_ServiceAccessToken true  "auths.HttpReq_ServiceAccessToken"
// @Success     200 {object} auths.HttpRsp_AccessTokenResponse
func (ctl ControlVanilla) Tenant(ctx echo.Context) (err error) {
	body := new(auths.HttpReq_ServiceAccessToken)
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
	tenant := tenants.Tenant{
		Hash: hash,
	}
	tenant_cond := stmt.Equal("hash", tenant.Hash)
	err = ctl.dialect.QueryRows(tenant.TableName(), []string{"id"}, tenant_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, i int) error {
			err := scan.Scan(&id)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get tenant")
	}

	time_now := time.Now()

	if id == 0 {
		// save tenant
		new_tenant := tenants.NewTenant(hash, pattern, time_now)
		updatecolumns := []string{"hash"}
		_, lastid, err := ctl.dialect.InsertOrUpdate(new_tenant.TableName(), new_tenant.ColumnNames(), updatecolumns, new_tenant.Values())(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "failed to save new tenant")
		}
		id = lastid // get tenant id
	}

	claims := auths.TenantAccessTokenClaims{
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

	return ctx.JSON(http.StatusOK, auths.HttpRsp_AccessTokenResponse{
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
