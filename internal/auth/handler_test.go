package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -------- Mock Service --------
type mockService struct {
	mock.Mock
}

func (m *mockService) Register(input RegisterRequest) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockService) Login(input LoginRequest) (*User, string, string, error) {
	args := m.Called(input)
	user, _ := args.Get(0).(*User)
	return user, args.String(1), args.String(2), args.Error(3)
}

func (m *mockService) RefreshToken(refreshToken string) (*User, string, error) {
	args := m.Called(refreshToken)
	user, _ := args.Get(0).(*User)
	return user, args.String(1), args.Error(2)
}

func (m *mockService) GetProfile(userID uuid.UUID) (*User, error) {
	args := m.Called(userID)
	user, _ := args.Get(0).(*User)
	return user, args.Error(1)
}

// -------- Helper --------
func performRequest(r http.Handler, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	r.ServeHTTP(w, req)
	return w
}

// -------- Init Logger --------
func init() {
	config.Log = logrus.New()
	config.Log.Out = io.Discard
}

// -------- Tests --------
func TestRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	reqBody := RegisterRequest{Name: "test", Email: "test@mail.com", Password: "secret"}
	mockSvc.On("Register", mock.Anything).Return(nil)

	r := gin.New()
	r.POST("/register", handler.Register)

	w := performRequest(r, "POST", "/register", reqBody)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
	mockSvc.AssertExpectations(t)
}

func TestRegister_Failed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	reqBody := RegisterRequest{Name: "test", Email: "fail@mail.com", Password: "secret"}
	mockSvc.On("Register", mock.Anything).Return(errors.New("failed"))

	r := gin.New()
	r.POST("/register", handler.Register)

	w := performRequest(r, "POST", "/register", reqBody)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed")
	mockSvc.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	reqBody := LoginRequest{Email: "test@mail.com", Password: "secret"}
	expectedUser := &User{Email: "test@mail.com"}

	mockSvc.On("Login", mock.Anything).Return(expectedUser, "access123", "refresh123", nil)

	r := gin.New()
	r.POST("/login", handler.Login)

	w := performRequest(r, "POST", "/login", reqBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test@mail.com")

	cookies := w.Result().Cookies()
	assert.Equal(t, "access123", cookies[0].Value)
	assert.Equal(t, "refresh123", cookies[1].Value)
	mockSvc.AssertExpectations(t)
}

func TestLogin_Failed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	reqBody := LoginRequest{Email: "wrong@mail.com", Password: "wrong"}
	mockSvc.On("Login", mock.Anything).Return(nil, "", "", errors.New("invalid credentials"))

	r := gin.New()
	r.POST("/login", handler.Login)

	w := performRequest(r, "POST", "/login", reqBody)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
	mockSvc.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	expectedUser := &User{Email: "ref@mail.com"}
	mockSvc.On("RefreshToken", "refresh123").Return(expectedUser, "newAccess123", nil)

	r := gin.New()
	r.GET("/refresh", handler.Refresh)

	cookie := &http.Cookie{Name: "refresh_token", Value: "refresh123"}
	w := performRequest(r, "GET", "/refresh", nil, cookie)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ref@mail.com")
	assert.Equal(t, "newAccess123", w.Result().Cookies()[0].Value)
	mockSvc.AssertExpectations(t)
}

func TestRefresh_Failed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	mockSvc.On("RefreshToken", "badtoken").Return(nil, "", errors.New("invalid token"))

	r := gin.New()
	r.GET("/refresh", handler.Refresh)

	cookie := &http.Cookie{Name: "refresh_token", Value: "badtoken"}
	w := performRequest(r, "GET", "/refresh", nil, cookie)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
	mockSvc.AssertExpectations(t)
}

func TestProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	userID := uuid.New()
	expectedUser := &User{ID: userID, Email: "profile@mail.com"}

	mockSvc.On("GetProfile", userID).Return(expectedUser, nil)

	r := gin.New()
	r.GET("/profile", func(c *gin.Context) {
		c.Set("user_id", userID.String())
		handler.Profile(c)
	})

	w := performRequest(r, "GET", "/profile", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "profile@mail.com")
	mockSvc.AssertExpectations(t)
}

func TestProfile_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockService)
	handler := NewHandler(mockSvc)

	r := gin.New()
	r.GET("/profile", handler.Profile)

	w := performRequest(r, "GET", "/profile", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(nil)

	r := gin.New()
	r.GET("/logout", handler.Logout)

	w := performRequest(r, "GET", "/logout", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logged out")

	// pastikan cookie dihapus
	found := false
	for _, c := range w.Result().Cookies() {
		if c.Name == "access_token" && c.Value == "" {
			found = true
		}
	}
	assert.True(t, found)
}