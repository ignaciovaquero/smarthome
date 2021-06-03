package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockContext interface {
	echo.Context
	GetToken(secret string) *jwt.Token
}

type baseMockContext struct {
	Token map[string]string
}

func (base *baseMockContext) GetToken(secret string) *jwt.Token {
	tok, _ := jwt.Parse(base.Token["token"], func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return tok
}

func (base *baseMockContext) Request() *http.Request {
	return nil
}

func (base *baseMockContext) SetRequest(r *http.Request) {
	return
}

func (base *baseMockContext) SetResponse(r *echo.Response) {
	return
}

func (base *baseMockContext) Response() *echo.Response {
	return nil
}

func (base *baseMockContext) IsTLS() bool {
	return false
}

func (base *baseMockContext) IsWebSocket() bool {
	return false
}

func (base *baseMockContext) Scheme() string {
	return ""
}

func (base *baseMockContext) RealIP() string {
	return ""
}

func (base *baseMockContext) Path() string {
	return ""
}

func (base *baseMockContext) SetPath(p string) {
	return
}

func (base *baseMockContext) Param(name string) string {
	return ""
}

func (base *baseMockContext) ParamNames() []string {
	return []string{}
}

func (base *baseMockContext) SetParamNames(names ...string) {
	return
}

func (base *baseMockContext) ParamValues() []string {
	return []string{}
}

func (base *baseMockContext) SetParamValues(values ...string) {
	return
}

func (base *baseMockContext) QueryParam(name string) string {
	return ""
}

func (base *baseMockContext) QueryParams() url.Values {
	return url.Values{}
}

func (base *baseMockContext) QueryString() string {
	return ""
}

func (base *baseMockContext) FormValue(name string) string {
	return ""
}

func (base *baseMockContext) FormParams() (url.Values, error) {
	return url.Values{}, nil
}

func (base *baseMockContext) FormFile(name string) (*multipart.FileHeader, error) {
	return nil, nil
}

func (base *baseMockContext) MultipartForm() (*multipart.Form, error) {
	return nil, nil
}

func (base *baseMockContext) Cookie(name string) (*http.Cookie, error) {
	return nil, nil
}

func (base *baseMockContext) SetCookie(cookie *http.Cookie) {
	return
}

func (base *baseMockContext) Cookies() []*http.Cookie {
	return []*http.Cookie{}
}

func (base *baseMockContext) Get(key string) interface{} {
	return nil
}

func (base *baseMockContext) Set(key string, val interface{}) {
	return
}

func (base *baseMockContext) Bind(i interface{}) error {
	return nil
}

func (base *baseMockContext) Validate(i interface{}) error {
	return nil
}

func (base *baseMockContext) Render(code int, name string, data interface{}) error {
	return nil
}

func (base *baseMockContext) HTML(code int, html string) error {
	return nil
}

func (base *baseMockContext) HTMLBlob(code int, b []byte) error {
	return nil
}

func (base *baseMockContext) String(code int, s string) error {
	return nil
}

func (base *baseMockContext) JSON(code int, i interface{}) error {
	if t, ok := i.(map[string]string); ok {
		base.Token = t
	}
	return nil
}

func (base *baseMockContext) JSONPretty(code int, i interface{}, indent string) error {
	return nil
}

func (base *baseMockContext) JSONBlob(code int, b []byte) error {
	return nil
}

func (base *baseMockContext) JSONP(code int, callback string, i interface{}) error {
	return nil
}

func (base *baseMockContext) JSONPBlob(code int, callback string, b []byte) error {
	return nil
}

func (base *baseMockContext) XML(code int, i interface{}) error {
	return nil
}

func (base *baseMockContext) XMLPretty(code int, i interface{}, indent string) error {
	return nil
}

func (base *baseMockContext) XMLBlob(code int, b []byte) error {
	return nil
}

func (base *baseMockContext) Blob(code int, contentType string, b []byte) error {
	return nil
}

func (base *baseMockContext) Stream(code int, contentType string, r io.Reader) error {
	return nil
}

func (base *baseMockContext) File(file string) error {
	return nil
}

func (base *baseMockContext) Attachment(file string, name string) error {
	return nil
}

func (base *baseMockContext) Inline(file string, name string) error {
	return nil
}

func (base *baseMockContext) NoContent(code int) error {
	return nil
}

func (base *baseMockContext) Redirect(code int, url string) error {
	return nil
}

func (base *baseMockContext) Error(err error) {
	return
}

func (base *baseMockContext) Handler() echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}

func (base *baseMockContext) SetHandler(h echo.HandlerFunc) {
	return
}

func (base *baseMockContext) Logger() echo.Logger {
	return nil
}

func (base *baseMockContext) SetLogger(l echo.Logger) {
	return
}

func (base *baseMockContext) Echo() *echo.Echo {
	return nil
}

func (base *baseMockContext) Reset(r *http.Request, w http.ResponseWriter) {
	return
}

type correctMockContext struct {
	baseMockContext
}

func (c *correctMockContext) Request() *http.Request {
	authParams := `{"username": "admin", "password": "admin"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(authParams)))
	return &http.Request{
		Body: r,
	}
}

type invalidPayloadMockContext struct {
	baseMockContext
}

func (i *invalidPayloadMockContext) Request() *http.Request {
	authParams := "Not a json"
	r := ioutil.NopCloser(bytes.NewReader([]byte(authParams)))
	return &http.Request{
		Body: r,
	}
}

type mockSmartHomeError struct {
	mockSmartHome
}

func (*mockSmartHomeError) Authenticate(username, password string) error {
	return fmt.Errorf("Error")
}

func (*mockSmartHomeError) SetCredentials(username, password string) error {
	return fmt.Errorf("Error")
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name              string
		ctx               mockContext
		cl                *Client
		httpErrorExpected bool
		errorExpected     bool
	}{
		{
			name: "Request with correct payload, no authentication errors",
			ctx:  &correctMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			httpErrorExpected: false,
			errorExpected:     false,
		},
		{
			name: "Request with invalid payload",
			ctx:  &invalidPayloadMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			httpErrorExpected: true,
			errorExpected:     false,
		},
		{
			name: "Request with correct payload, authentication error",
			ctx:  &correctMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHomeError{},
			),
			httpErrorExpected: false,
			errorExpected:     true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.Login(tc.ctx)
			if tc.httpErrorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			if tc.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			tok := tc.ctx.GetToken(tc.cl.Config.JWTSecret)
			assert.True(tt, tok.Valid)
		})
	}
}

func TestSignUp(t *testing.T) {
	testCases := []struct {
		name              string
		ctx               mockContext
		cl                *Client
		httpErrorExpected bool
		errorExpected     bool
	}{
		{
			name: "Request with correct payload, no errors",
			ctx:  &correctMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			httpErrorExpected: false,
			errorExpected:     false,
		},
		{
			name: "Request with invalid payload",
			ctx:  &invalidPayloadMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			httpErrorExpected: true,
			errorExpected:     false,
		},
		{
			name: "Request with correct payload, controller error",
			ctx:  &correctMockContext{},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHomeError{},
			),
			httpErrorExpected: false,
			errorExpected:     true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.SignUp(tc.ctx)
			if tc.httpErrorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			if tc.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
