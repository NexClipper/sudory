package control

import (
	"reflect"
	"runtime"
)

// type Context interface {
// 	//echo
// 	Echo() echo.Context
// 	SetEcho(e echo.Context) Context
// 	Forms() map[string]string
// 	FormString() string
// 	Params() map[string]string
// 	ParamString() string
// 	Queries() map[string]string
// 	QueryString() string
// 	Body() []byte
// 	Bind(interface{}) error
// 	Object() interface{}
// 	//database
// 	Database() database.Context
// 	SetDatabase(database.Context) Context
// }

// type OnceMapStringString struct {
// 	sync.Once
// 	v map[string]string
// }

// type OnceBytes struct {
// 	sync.Once
// 	v []byte
// }

// type RequestValue struct {
// 	//echo
// 	echo                    echo.Context
// 	param, query, formParam OnceMapStringString
// 	body                    OnceBytes
// 	object                  interface{}
// 	//database
// 	db database.Context
// }

// // func (holder RequestValue) TicketId() uint64 {
// // 	return holder.ticketId
// // }
// // func (holder RequestValue) SetTicketId(t uint64) Contexter {
// // 	holder.ticketId = t
// // 	return &holder
// // }

// func (holder RequestValue) Echo() echo.Context {
// 	return holder.echo
// }
// func (holder RequestValue) SetEcho(e echo.Context) Context {
// 	holder.echo = e
// 	return &holder
// }

// func (holder *RequestValue) Params() map[string]string {
// 	holder.param.Do(func() {
// 		holder.param.v = make(map[string]string)
// 		for _, name := range holder.echo.ParamNames() {
// 			holder.param.v[name] = holder.echo.Param(name)
// 		}
// 	})
// 	return holder.param.v
// }
// func (holder *RequestValue) ParamString() string {
// 	s := make([]string, 0, len(holder.Params()))
// 	for key := range holder.Params() {
// 		s = append(s, fmt.Sprintf("%s:%s", key, holder.Params()[key]))
// 	}
// 	return strings.Join(s, ",")
// }

// func (holder *RequestValue) Queries() map[string]string {
// 	holder.query.Do(func() {
// 		holder.query.v = make(map[string]string)
// 		for key := range holder.echo.QueryParams() {
// 			holder.query.v[key] = holder.echo.QueryParam(key)
// 		}
// 	})
// 	return holder.query.v
// }

// func (holder *RequestValue) QueryString() string {
// 	return holder.echo.QueryString()
// }

// func (holder *RequestValue) Forms() map[string]string {
// 	holder.formParam.Do(func() {
// 		holder.formParam.v = make(map[string]string)
// 		formdatas, err := holder.echo.FormParams()
// 		if err != nil {
// 			return
// 		}
// 		for key := range formdatas {
// 			holder.formParam.v[key] = holder.echo.FormValue(key)
// 		}
// 	})
// 	return holder.formParam.v
// }
// func (holder *RequestValue) FormString() string {
// 	s := make([]string, 0, len(holder.Forms()))
// 	for key := range holder.Forms() {
// 		s = append(s, fmt.Sprintf("%s=%s", key, holder.Forms()[key]))
// 	}
// 	return strings.Join(s, "&")
// }

// func (holder *RequestValue) Body() []byte {
// 	holder.body.Do(func() {
// 		//body read all
// 		//ranout buffer
// 		holder.body.v, _ = ioutil.ReadAll(holder.echo.Request().Body) //read all body
// 		//restore
// 		holder.echo.Request().Body = ioutil.NopCloser(bytes.NewBuffer(holder.body.v))
// 	})
// 	return holder.body.v
// }

// func (holder *RequestValue) Bind(v interface{}) error {

// 	if err := holder.echo.Bind(v); err != nil {
// 		return err
// 	}

// 	// if err := json.Unmarshal(holder.Body(), v); err != nil {
// 	// 	return err
// 	// }
// 	holder.object = v
// 	return nil
// }

// func (holder *RequestValue) Object() interface{} {
// 	return holder.object
// }

// func (holder RequestValue) Database() database.Context {
// 	return holder.db
// }
// func (holder RequestValue) SetDatabase(d database.Context) Context {
// 	holder.db = d
// 	return &holder
// }

// func HttpJsonResponsor(ctx echo.Context, status int, v interface{}) error {
// 	return ctx.JSON(status, v)
// }

func OK() interface{} {
	return "OK"
}

// func Lock(engine *xorm.Engine, ctx Context, operate func(Context) (interface{}, error)) (interface{}, error) {
// 	// return func(ctx Context, operate func(Context) (interface{}, error)) (interface{}, error) {
// 	return engine.Transaction(func(s *xorm.Session) (interface{}, error) {
// 		ctx = ctx.SetDatabase(database.NewXormContext(s)) //new database context
// 		return operate(ctx)
// 	})
// 	// }
// }

// func LOCK(engine *xorm.Engine, operate func(*xorm.Session) error) error {
// 	_, err := engine.Transaction(func(s *xorm.Session) (interface{}, error) {
// 		return nil, operate(s)
// 	})

// 	return err
// }

// func Compose(engine *xorm.Engine, operator HandlerFunc, validator ...ValidatorFunc) echo.HandlerFunc {
// 	var (
// 		context Context = &RequestValue{}
// 	)

// 	exec_validator := func(ctx Context, validator ValidatorFunc) (int, error) {
// 		var (
// 			code int
// 			err  error
// 		)
// 		block.Block{
// 			Try: func() {
// 				_, lockerr := Lock(engine, context, func(ctx Context) (interface{}, error) {
// 					code, err = validator(context)
// 					return code, err
// 				})
// 				if err == nil && lockerr != nil {
// 					err = errors.Wrapf(lockerr, "xorm commit")
// 				}
// 			},
// 			Catch: func(ex error) {
// 				err = errors.Wrapf(ex, "catch")
// 			},
// 		}.Do()

// 		return code, err
// 	}

// 	exec_operator := func(ctx Context, operator HandlerFunc) (interface{}, error) {
// 		var (
// 			v   interface{}
// 			err error
// 		)
// 		block.Block{
// 			Try: func() {
// 				_, lockerr := Lock(engine, context, func(ctx Context) (interface{}, error) {
// 					v, err = operator(context)
// 					return v, err
// 				})
// 				if err == nil && lockerr != nil {
// 					err = errors.Wrapf(lockerr, "xorm commit")
// 				}
// 			},
// 			Catch: func(ex error) {
// 				err = errors.Wrapf(ex, "catch")
// 			},
// 		}.Do()

// 		return v, err
// 	}

// 	return func(echo_ echo.Context) error {
// 		context = context.SetEcho(echo_)

// 		//pre exec
// 		for _, validator := range validator {
// 			code, err := exec_validator(context, validator)
// 			if err != nil {
// 				err = errors.Wrapf(err, "exec %s", validator.Name())

// 				if code == http.StatusOK {
// 					code = http.StatusBadRequest
// 				}

// 				echo_.String(code, err.Error())
// 				return err
// 			}
// 		}

// 		//main exec
// 		v, err := exec_operator(context, operator)
// 		if err != nil {
// 			err = errors.Wrapf(err, "exec %s", operator.Name())
// 			echo_.String(http.StatusInternalServerError, err.Error())
// 			return err
// 		}

// 		echo_.JSON(http.StatusOK, v)
// 		return nil
// 	}
// }

// func WrapValidator(fn func(e echo.Context) error) func(e echo.Context) (interface{}, error) {
// 	return func(e echo.Context) (interface{}, error) {
// 		return nil, fn(e)
// 	}
// }

// func GetFunctionName(i interface{}) string {
// 	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
// }

// type HandlerFunc func(Context) (interface{}, error)

// func (h HandlerFunc) Name() string {
// 	return TypeName(h)
// }

// type ValidatorFunc func(Context) (int, error)

// func (h ValidatorFunc) Name() string {
// 	return TypeName(h)
// }

func TypeName(i interface{}) string {
	t := reflect.ValueOf(i).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	}
	return t.String()
}

// func Param(echo_ echo.Context) map[string]string {
// 	m := map[string]string{}
// 	for _, name := range echo_.ParamNames() {
// 		m[name] = echo_.Param(name)
// 	}

// 	return m
// }

// func ParamString(echo_ echo.Context) string {
// 	param := Param(echo_)
// 	s := make([]string, 0, len(param))
// 	for key := range param {
// 		s = append(s, fmt.Sprintf("%s:%s", key, param[key]))
// 	}
// 	return strings.Join(s, ",")
// }

// func QueryParam(echo_ echo.Context) map[string]string {
// 	m := map[string]string{}
// 	for key := range echo_.QueryParams() {
// 		m[key] = echo_.QueryParam(key)
// 	}

// 	return m
// }

// func QueryParamString(echo_ echo.Context) string {
// 	return echo_.QueryString()
// }

// func FormParam(echo_ echo.Context) map[string]string {
// 	m := map[string]string{}
// 	formdatas, err := echo_.FormParams()
// 	if err != nil {
// 		return m
// 	}

// 	for key := range formdatas {
// 		m[key] = echo_.FormValue(key)
// 	}

// 	return m
// }

// func FormParamString(echo_ echo.Context) string {
// 	formparam := FormParam(echo_)
// 	s := make([]string, 0, len(formparam))
// 	for key := range formparam {
// 		s = append(s, fmt.Sprintf("%s=%s", key, formparam[key]))
// 	}
// 	return strings.Join(s, "&")
// }

// func Body(echo_ echo.Context) []byte {
// 	//body read all
// 	b, _ := ioutil.ReadAll(echo_.Request().Body) //read all body //ranout
// 	//restore
// 	echo_.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

// 	return b
// }

// func Bind(echo_ echo.Context, v interface{}) error {
// 	if err := echo_.Bind(v); err != nil {
// 		return err
// 	}
// 	return nil
// }
