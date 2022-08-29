package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Tags to construct the logger format.
		//
		// - time_unix
		// - time_unix_nano
		// - time_rfc3339
		// - time_rfc3339_nano
		// - time_custom
		// - id (Request ID)
		// - remote_ip
		// - uri
		// - host
		// - method
		// - path
		// - protocol
		// - referer
		// - user_agent
		// - status
		// - error
		// - latency (In nanoseconds)
		// - latency_human (Human readable)
		// - bytes_in (Bytes received)
		// - bytes_out (Bytes sent)
		// - header:<NAME>
		// - query:<NAME>
		// - form:<NAME>
		//
		// Example "${remote_ip} ${status}"
		//
		// Optional. Default value DefaultLoggerConfig.Format.
		Format string `yaml:"format"`

		// Optional. Default value DefaultLoggerConfig.CustomTimeFormat.
		CustomTimeFormat string `yaml:"custom_time_format"`

		// Output is a writer where logs in JSON format are written.
		// Optional. Default value os.Stdout.
		Output io.Writer

		// template *fasttemplate.Template
		colorer *color.Color
		pool    *sync.Pool
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		colorer:          color.New(),
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Format == "" {
		config.Format = DefaultLoggerConfig.Format
	}
	if config.Output == nil {
		config.Output = DefaultLoggerConfig.Output
	}

	texts, tags, err := ParseTemplate(config.Format, "${", "}")
	if err != nil {
		panic(err)
	}
	// config.template = fasttemplate.New(config.Format, "${", "}")
	config.colorer = color.New()
	config.colorer.SetOutput(config.Output)
	config.pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			buf := config.pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer config.pool.Put(buf)

			// if _, err = config.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			tagBuilder := func(tag string) (string, bool) {
				switch tag {
				case "time_unix":
					return strconv.FormatInt(time.Now().Unix(), 10), true
				case "time_unix_nano":
					return strconv.FormatInt(time.Now().UnixNano(), 10), true
				case "time_rfc3339":
					return time.Now().Format(time.RFC3339), true
				case "time_rfc3339_nano":
					return time.Now().Format(time.RFC3339Nano), true
				case "time_custom":
					return time.Now().Format(config.CustomTimeFormat), true
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					return id, true
				case "remote_ip":
					return c.RealIP(), true
				case "host":
					return req.Host, true
				case "uri":
					return req.RequestURI, true
				case "method":
					return req.Method, true
				case "path":
					p := req.URL.Path
					if p == "" {
						p = "/"
					}
					return p, true
				case "protocol":
					return req.Proto, true
				case "referer":
					return req.Referer(), true
				case "user_agent":
					return req.UserAgent(), true
				case "status":
					n := res.Status
					s := config.colorer.Green(n)
					switch {
					case n >= 500:
						s = config.colorer.Red(n)
					case n >= 400:
						s = config.colorer.Yellow(n)
					case n >= 300:
						s = config.colorer.Cyan(n)
					}
					return s, true
				case "error":
					if err != nil {
						// Error may contain invalid JSON e.g. `"`
						b, _ := json.Marshal(err.Error())
						b = b[1 : len(b)-1]
						return string(b), true
					}
				case "latency":
					l := stop.Sub(start)
					return strconv.FormatInt(int64(l), 10), true
				case "latency_human":
					return stop.Sub(start).String(), true
				case "bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					return cl, true
				case "bytes_out":
					return strconv.FormatInt(res.Size, 10), true
				default:
					switch {
					case strings.HasPrefix(tag, "header:"):
						return c.Request().Header.Get(tag[7:]), true
					case strings.HasPrefix(tag, "query:"):
						return c.QueryParam(tag[6:]), true
					case strings.HasPrefix(tag, "form:"):
						return c.FormValue(tag[5:]), true
					case strings.HasPrefix(tag, "cookie:"):
						cookie, err := c.Cookie(tag[7:])
						if err == nil {
							return (cookie.Value), true
						}
					}
				}
				return "", false
			}
			// return 0, nil
			// });

			_, err = func() (nn int, err error) {
				var n int
				for i := range tags {
					value, ok := tagBuilder(tags[i])
					if ok {
						n, err = buf.Write(texts[i])
						nn += n
						if err != nil {
							return
						}
						n, err = buf.WriteString(value)
						nn += n
						if err != nil {
							return
						}
					}
				}
				n, err = buf.Write(texts[len(texts)-1])
				nn += n
				if err != nil {
					return
				}
				return
			}()
			if err != nil {
				return
			}

			if config.Output == nil {
				_, err = c.Logger().Output().Write(buf.Bytes())
				return
			}
			_, err = config.Output.Write(buf.Bytes())
			return
		}
	}
}

func ParseTemplate(template, starttag, endtag string) (texts [][]byte, tags []string, err error) {
	s, a, b := unsafeString2Bytes(template), unsafeString2Bytes(starttag), unsafeString2Bytes(endtag)

	tagsCount := bytes.Count(s, a)
	if tagsCount == 0 {
		return
	}

	texts = make([][]byte, 0, tagsCount+1)
	tags = make([]string, 0, tagsCount)

	for {
		n := bytes.Index(s, a)
		if n < 0 {
			texts = append(texts, s)
			break
		}
		texts = append(texts, s[:n])

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 {
			err = fmt.Errorf("cannot find end tag=%q in the template=%q starting from %q", endtag, template, s)
			return
		}

		tags = append(tags, unsafeBytes2String(s[:n]))
		s = s[n+len(b):]
	}

	return
}

func unsafeBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeString2Bytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
