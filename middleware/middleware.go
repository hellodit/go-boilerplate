package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Middleware struct {}

func Init() *Middleware {
	return &Middleware{}
}

func makeLogEntry(c echo.Context) *logrus.Entry {
	if c == nil {
		return logrus.WithFields(logrus.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return logrus.WithFields(logrus.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri": c.Request().URL.String(),
		"ip": c.Request().RemoteAddr,
	})
}

func (m *Middleware) MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("Incoming request")
		return next(c)
	}
}

func (m *Middleware) ErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if report.Internal == nil{
		report.SetInternal(errors.New(""))
	}


	makeLogEntry(c).Error(report.Message)
	c.JSON(report.Code, map[string]map[string]interface{}{
		"error": {
			"code": report.Code,
			"message": report.Internal.Error(),
			"errors":  report.Message,
		},
	})
}



func (m *Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if !strings.Contains(tokenString, "Bearer") {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Token not provided"))
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, errors.New("invalid token")
			}

			return []byte(viper.GetString("JWT_SECRET")), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}
}


// Logrus : implement Logger
type Logrus struct {
	*logrus.Logger
}

// Logger ...
var Logger *logrus.Logger

// GetEchoLogger for e.Logger
func GetEchoLogger() Logrus {
	return Logrus{Logger}
}

// Level returns logger level
func (l Logrus) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Panic("Invalid level")
	}

	return log.OFF
}

// SetHeader is a stub to satisfy interface
// It's controlled by Logger
func (l Logrus) SetHeader(_ string) {}

// SetPrefix It's controlled by Logger
func (l Logrus) SetPrefix(s string) {}

// Prefix It's controlled by Logger
func (l Logrus) Prefix() string {
	return ""
}

// SetLevel set level to logger from given log.Lvl
func (l Logrus) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Panic("Invalid level")
	}
}

// Output logger output func
func (l Logrus) Output() io.Writer {
	return l.Out
}

// SetOutput change output, default os.Stdout
func (l Logrus) SetOutput(w io.Writer) {
	Logger.SetOutput(w)
}

// Printj print json log
func (l Logrus) Printj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Print()
}

// Debugj debug json log
func (l Logrus) Debugj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Debug()
}

// Infoj info json log
func (l Logrus) Infoj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Info()
}

// Warnj warning json log
func (l Logrus) Warnj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Warn()
}

// Errorj error json log
func (l Logrus) Errorj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Error()
}

// Fatalj fatal json log
func (l Logrus) Fatalj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Fatal()
}

// Panicj panic json log
func (l Logrus) Panicj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Panic()
}

// Print string log
func (l Logrus) Print(i ...interface{}) {
	Logger.Print(i[0].(string))
}

// Debug string log
func (l Logrus) Debug(i ...interface{}) {
	Logger.Debug(i[0].(string))
}

// Info string log
func (l Logrus) Info(i ...interface{}) {
	Logger.Info(i[0].(string))
}

// Warn string log
func (l Logrus) Warn(i ...interface{}) {
	Logger.Warn(i[0].(string))
}

// Error string log
func (l Logrus) Error(i ...interface{}) {
	Logger.Error(i[0].(string))
}

// Fatal string log
func (l Logrus) Fatal(i ...interface{}) {
	Logger.Fatal(i[0].(string))
}

// Panic string log
func (l Logrus) Panic(i ...interface{}) {
	Logger.Panic(i[0].(string))
}

func logrusMiddlewareHandler(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()
	start := time.Now()
	if err := next(c); err != nil {
		c.Error(err)
	}
	stop := time.Now()

	p := req.URL.Path

	bytesIn := req.Header.Get(echo.HeaderContentLength)

	Logger.WithFields(map[string]interface{}{
		"time_rfc3339":  time.Now().Format(time.RFC3339),
		"remote_ip":     c.RealIP(),
		"host":          req.Host,
		"uri":           req.RequestURI,
		"method":        req.Method,
		"path":          p,
		"referer":       req.Referer(),
		"user_agent":    req.UserAgent(),
		"status":        res.Status,
		"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
		"latency_human": stop.Sub(start).String(),
		"bytes_in":      bytesIn,
		"bytes_out":     strconv.FormatInt(res.Size, 10),
	}).Infoln("Handled request")

	return nil
}

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return logrusMiddlewareHandler(c, next)
	}
}

// Hook is a function to process middleware.
func Hook() echo.MiddlewareFunc {
	return logger
}