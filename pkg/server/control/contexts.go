package control

import (
	"context"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	auth "github.com/NexClipper/sudory/pkg/server/model/auth/v2"
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

func GetClientSessionClaims(ctx echo.Context, tx stmtex.Preparer, dialect string) (*sessionv3.ClientSessionPayload, error) {
	const (
		// __HTTP_HEADER_KEY__   = __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__
		__HTTP_HEADER_KEY__   = "x-sudory-client-token"
		__CONTEXT_VALUE_KEY__ = ClientSession
	)

	switch v := GetContextValue(ctx.Request().Context(), __CONTEXT_VALUE_KEY__).(type) {
	case error:
		return nil, v
	case *sessionv3.ClientSessionPayload:
		return v, nil
	default:
		return func() (claims *sessionv3.ClientSessionPayload, err error) {
			time_now := time.Now()
			// get session token
			header_string := ctx.Request().Header.Get(__HTTP_HEADER_KEY__)

			if len(header_string) == 0 {
				err = errors.Errorf("missing request header%v", logs.KVL(
					"header", __HTTP_HEADER_KEY__,
				))
			}
			if err != nil {
				return
			}

			// parse to payload
			claims = new(sessionv3.ClientSessionPayload)
			token, err := jwt.ParseWithClaims(header_string, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ClientSession.SignatureSecret()), nil
			})
			err = errors.Wrapf(err, "jwt parse claims")
			if err != nil {
				return
			}

			// check type cast and valid
			if _, ok := token.Claims.(*sessionv3.ClientSessionPayload); !ok || !token.Valid {
				err = errors.Errorf("jwt verify%v", logs.KVL(
					"token", header_string,
				))
			}
			if err != nil {
				return
			}

			// smart polling interval
			// update echo request context
			echo_context := ctx.Request().Context()
			echo_context = SetContextValue(echo_context, __CONTEXT_VALUE_KEY__, claims)
			ctx.SetRequest(ctx.Request().WithContext(echo_context))

			// refresh session token
			// get cluster from db
			cluster := clusterv3.Cluster{}
			cluster.Uuid = claims.ClusterUuid
			cluster_cond := stmt.And(
				stmt.Equal("uuid", cluster.Uuid),
				stmt.IsNull("deleted"),
			)
			err = stmtex.Select(cluster.TableName(), cluster.ColumnNames(), cluster_cond, nil, nil).
				QueryRowContext(ctx.Request().Context(), tx, dialect)(func(s stmtex.Scanner) (err error) {
				err = cluster.Scan(s)
				err = errors.Wrapf(err, "cluster Scan")
				return
			})
			err = errors.Wrapf(err, "failed to get cluster")
			if err != nil {
				return
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
			func() (err error) {
				cluster_info := clusterinfov2.ClusterInformation{}
				columnnames := []string{"polling_count"}
				cond := stmt.Equal("cluster_uuid", claims.ClusterUuid)

				err = stmtex.Select(cluster_info.TableName(), columnnames, cond, nil, nil).
					QueryRowsContext(ctx.Request().Context(), tx, dialect)(func(scan stmtex.Scanner, _ int) error {
					return scan.Scan(&polling_count)
				})
				if err != nil {
					return errors.Wrapf(err, "failed to get service offset")
				}
				return
			}()

			//reflesh claims
			polling_interval := clusterv3.ConvPollingOption(cluster.PollingOption).Interval(time.Duration(int64(globvar.ClientConfig.PollInterval())*int64(time.Second)), polling_count) / time.Second
			claims.PollInterval = int(polling_interval)
			claims.ExpiresAt = globvar.ClientSession.ExpirationTime(time_now).Unix()
			claims.Loglevel = globvar.ClientConfig.Loglevel()

			// make response session token
			new_token, err := jwt.NewWithClaims(usedJwtSigningMethod(*token, jwt.SigningMethodHS256), claims).
				SignedString([]byte(globvar.ClientSession.SignatureSecret()))
			err = errors.Wrapf(err, "failed to make client session token to formed jwt")
			if err != nil {
				return
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
			session_found, err := stmtex.ExistContext(session.TableName(), session_cond)(ctx.Request().Context(), tx, dialect)
			err = errors.Wrapf(err, "failed to found client session")
			if err != nil {
				return
			}

			if !session_found {
				err = errors.Errorf("not found client session%v", logs.KVL(
					"uuid", session.Uuid,
				))
			}
			if err != nil {
				return
			}

			// update refreshed client session
			keys_values := map[string]interface{}{
				"token":           session.Token,
				"expiration_time": session.ExpirationTime,
				"updated":         session.Updated,
			}

			_, err = stmtex.Update(session.TableName(), keys_values, session_cond).
				ExecContext(ctx.Request().Context(), tx, dialect)
			err = errors.Wrapf(err, "failed to refreshed client session%v", logs.KVL(
				"uuid", claims.Uuid,
				"data", keys_values,
			))

			if err != nil {
				return
			}

			return GetClientSessionClaims(ctx, tx, dialect) // recurse
		}()
	}
}

func GetServiceAuthorizationClaims(ctx echo.Context) (*auth.ServiceAccessTokenClaims, error) {
	const (
		__CONTEXT_VALUE_KEY__ = ServiceClaims
	)

	switch v := GetContextValue(ctx.Request().Context(), __CONTEXT_VALUE_KEY__).(type) {
	case error:
		return nil, v
	case *auth.ServiceAccessTokenClaims:
		return v, nil
	default:
		return func() (claims *auth.ServiceAccessTokenClaims, err error) {
			// get bearer token
			bearer, token, ok := echoutil.ParseAuthorizationHeader(ctx.Request().Header)
			if !ok {
				err = errors.New("token is not available")
				return
			}
			if bearer != echoutil.AuthSchemaBearer {
				err = errors.New("token is not available")
				return
			}

			// parse jwt
			claims = new(auth.ServiceAccessTokenClaims)
			jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ServiceSession.SignatureSecret()), nil
			})
			if err != nil {
				return
			}
			// verify claims
			if _, ok := jwt_token.Claims.(*auth.ServiceAccessTokenClaims); !ok || !jwt_token.Valid {
				err = errors.Errorf("verify claims%v", logs.KVL(
					"token", token,
					"alg", jwt_token.Method.Alg(),
				))
			}
			if err != nil {
				return
			}

			// update echo request context
			echo_context := ctx.Request().Context()
			echo_context = SetContextValue(echo_context, __CONTEXT_VALUE_KEY__, claims)
			ctx.SetRequest(ctx.Request().WithContext(echo_context))

			return GetServiceAuthorizationClaims(ctx)
		}()
	}
}
