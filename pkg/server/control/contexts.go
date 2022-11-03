package control

import (
	"context"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/model/auths/v2"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	clusterinfov2 "github.com/NexClipper/sudory/pkg/server/model/cluster_infomation/v2"
	sessionv3 "github.com/NexClipper/sudory/pkg/server/model/session/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	ClientSession ContextValueKey = iota
	ServiceClaims
)

type ContextValueKey int

func GetContextValue(ctx context.Context, key ContextValueKey) interface{} {
	return ctx.Value(key)
}

func SetContextValue(ctx context.Context, key ContextValueKey, v interface{}) context.Context {
	return context.WithValue(ctx, key, v)
}

func GetClientSessionClaims(ctx echo.Context, tx excute.Preparer, dialect excute.SqlExcutor) (*sessionv3.ClientSessionPayload, error) {
	const (
		__HTTP_HEADER_KEY__   = "x-sudory-client-token"
		__CONTEXT_VALUE_KEY__ = ClientSession
	)

	var timeout = context.Background()

	switch v := GetContextValue(ctx.Request().Context(), __CONTEXT_VALUE_KEY__).(type) {
	case error:
		return nil, v
	case *sessionv3.ClientSessionPayload:
		return v, nil
	default:
		fn := func() interface{} {
			time_now := time.Now()
			// get session token
			header_string := ctx.Request().Header.Get(__HTTP_HEADER_KEY__)

			if len(header_string) == 0 {
				err := errors.Errorf("missing request header%v", logs.KVL(
					"header", __HTTP_HEADER_KEY__,
				))
				return err
			}

			// parse to payload
			claims := new(sessionv3.ClientSessionPayload)
			token, err := jwt.ParseWithClaims(header_string, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ClientSession.SignatureSecret()), nil
			})
			if err != nil {
				err = errors.Wrapf(err, "jwt parse claims")
				return err
			}

			// check type cast and valid
			if _, ok := token.Claims.(*sessionv3.ClientSessionPayload); !ok || !token.Valid {
				err = errors.Errorf("jwt verify%v", logs.KVL(
					"token", header_string,
				))
				return err
			}

			// refresh session token
			// get cluster from db
			cluster := clusterv3.Cluster{}
			cluster.Uuid = claims.ClusterUuid
			cluster_cond := stmt.And(
				stmt.Equal("uuid", cluster.Uuid),
				stmt.IsNull("deleted"),
			)
			err = dialect.QueryRow(cluster.TableName(), cluster.ColumnNames(), cluster_cond, nil, nil)(
				timeout, tx)(
				func(s excute.Scanner) error {
					err = cluster.Scan(s)
					err = errors.WithStack(err)

					return err
				})
			if err != nil {
				err = errors.Wrapf(err, "failed to get cluster")
				return err
			}

			// // count service from db
			// service := servicev2.Service_status{}
			// service.ClusterUuid = claims.ClusterUuid
			// service_cond := pollingServiceCondition(service.ClusterUuid)
			// service_count, err := vanilla.Stmt.Count(service.TableName(), service_cond, nil)(ctx.Request().Context(), tx)
			// err = errors.Wrapf(err, "failed to count service")
			// if err != nil {
			// 	return
			// }

			var polling_count int
			clusterinfo := clusterinfov2.ClusterInformation{}
			clusterinfo_columns := []string{"polling_count"}
			clusterinfo_cond := stmt.Equal("cluster_uuid", claims.ClusterUuid)

			err = dialect.QueryRows(clusterinfo.TableName(), clusterinfo_columns, clusterinfo_cond, nil, nil)(
				timeout, tx)(
				func(scan excute.Scanner, _ int) error {
					err := scan.Scan(&polling_count)
					err = errors.WithStack(err)

					return err
				})
			if err != nil {
				err = errors.Wrapf(err, "failed to get service offset")
				return err
			}

			//reflesh claims
			polling_interval := clusterv3.ConvPollingOption(cluster.PollingOption).Interval(time.Duration(int64(globvar.ClientConfig.PollInterval())*int64(time.Second)), polling_count) / time.Second
			claims.PollInterval = int(polling_interval)
			claims.ExpiresAt = globvar.ClientSession.ExpirationTime(time_now).Unix()
			claims.Loglevel = globvar.ClientConfig.Loglevel()

			// make response session token
			new_token, err := jwt.NewWithClaims(usedJwtSigningMethod(*token, jwt.SigningMethodHS256), claims).
				SignedString([]byte(globvar.ClientSession.SignatureSecret()))
			if err != nil {
				err = errors.Wrapf(err, "failed to make client session token to formed jwt")
				return err
			}

			// set response head
			ctx.Response().Header().Set(__HTTP_HEADER_KEY__, new_token)

			// update session token
			session := sessionv3.Session{}
			session.Uuid = claims.Uuid
			session.Token = new_token
			session.ExpirationTime = *vanilla.NewNullTime(time.Unix(claims.ExpiresAt, 0))
			session.Updated = *vanilla.NewNullTime(time_now)

			session_cond := stmt.And(
				stmt.Equal("uuid", session.Uuid),
				stmt.IsNull("deleted"),
			)

			// found session
			session_found, err := dialect.Exist(session.TableName(), session_cond)(timeout, tx)
			if err != nil {
				err = errors.Wrapf(err, "failed to found client session")
				return err
			}

			if !session_found {
				err = errors.Errorf("not found client session%v", logs.KVL(
					"uuid", session.Uuid,
				))
				return err
			}

			// update refreshed client session
			keys_values := map[string]interface{}{
				"token":           session.Token,
				"expiration_time": session.ExpirationTime,
				"updated":         session.Updated,
			}

			_, err = dialect.Update(session.TableName(), keys_values, session_cond)(timeout, tx)
			if err != nil {
				err = errors.Wrapf(err, "failed to refreshed client session%v", logs.KVL(
					"uuid", claims.Uuid,
					"data", keys_values,
				))
				return err
			}

			return claims
		}

		// smart polling interval
		// update echo request context
		echo_context := ctx.Request().Context()
		echo_context = SetContextValue(echo_context, __CONTEXT_VALUE_KEY__, fn())
		ctx.SetRequest(ctx.Request().WithContext(echo_context))

		return GetClientSessionClaims(ctx, tx, dialect) // recurse
	}
}

func GetServiceAuthorizationClaims(ctx echo.Context) (*auths.TenantAccessTokenClaims, error) {
	const (
		__CONTEXT_VALUE_KEY__  = ServiceClaims
		default_tenant_hash    = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
		default_tenant_pattern = ""
	)

	switch v := GetContextValue(ctx.Request().Context(), __CONTEXT_VALUE_KEY__).(type) {
	case error:
		return nil, v
	case *auths.TenantAccessTokenClaims:
		return v, nil
	default:
		fn := func() interface{} {
			// get http authorization header
			if len(echoutil.GetAuthorizationHeader(ctx.Request().Header)) == 0 {
				// treat the no tenant information as a default tenant
				claims := new(auths.TenantAccessTokenClaims)
				claims.ID = 1
				claims.Hash = default_tenant_hash
				claims.Tenant = default_tenant_pattern
				claims.IssuedAt = time.Now().Unix()
				claims.ExpiresAt = 0

				return claims
			}

			// parse bearer token
			bearer, token, ok := echoutil.ParseAuthorizationHeader(ctx.Request().Header)
			if !ok {
				err := errors.New("token is not available")
				return err
			}
			if bearer != echoutil.HTTP_AUTH_SCHEMA_BEARER {
				err := errors.New("token is not available")
				return err
			}

			// parse jwt
			claims := new(auths.TenantAccessTokenClaims)
			jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ServiceSession.SignatureSecret()), nil
			})
			if err != nil {
				return err
			}
			// verify claims
			if _, ok := jwt_token.Claims.(*auths.TenantAccessTokenClaims); !ok || !jwt_token.Valid {
				err = errors.Errorf("verify claims%v", logs.KVL(
					"token", token,
					"alg", jwt_token.Method.Alg(),
				))
				if err != nil {
					return err
				}
			}

			return claims
		}

		// update echo request context
		echo_context := ctx.Request().Context()
		echo_context = SetContextValue(echo_context, __CONTEXT_VALUE_KEY__, fn())
		ctx.SetRequest(ctx.Request().WithContext(echo_context))

		return GetServiceAuthorizationClaims(ctx)
	}
}
