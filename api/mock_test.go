package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type mockContext interface {
	echo.Context
	GetToken(secret string) *jwt.Token
}

type baseMockContext struct {
	Body        string
	Parameter   string
	JSONPayload interface{}
}

func (base *baseMockContext) GetToken(secret string) *jwt.Token {
	payload, ok := base.JSONPayload.(map[string]string)
	if !ok {
		return nil
	}
	t, ok := payload["token"]
	if !ok {
		return nil
	}
	tok, _ := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return tok
}

func (base *baseMockContext) Request() *http.Request {
	return &http.Request{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(base.Body))),
	}
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
	return base.Parameter
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
	base.JSONPayload = i
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

type mockSmartHome struct {
	Err error
}

func (m *mockSmartHome) Authenticate(username, password string) error {
	return m.Err
}
func (m *mockSmartHome) SetCredentials(username, password string) error {
	return m.Err
}
func (m *mockSmartHome) SetRoomOptions(room string, enabled bool, thresholdOn, thresholdOff float32) error {
	return m.Err
}
func (m *mockSmartHome) GetRoomOptions(room string) (map[string]types.AttributeValue, error) {
	return map[string]types.AttributeValue{}, nil
}
func (m *mockSmartHome) DeleteRoomOptions(room string) error {
	return m.Err
}
func (m *mockSmartHome) DeleteUser(username string) error {
	return m.Err
}
