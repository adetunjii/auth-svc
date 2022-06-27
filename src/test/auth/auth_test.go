package auth

import (
	"checklos/src/helpers"
	"checklos/src/services"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	router *gin.Engine
	mock   *service.MockUserService
}

func beforeEach(t *testing.T) Test {
	ctrl := gomock.NewController(t)
	mockUser := service.NewMockUserService(ctrl)

	router := gin.Default()
	h := &UserProto.Handler{
		UserService: mockUser,
	}
	services.DefineRouter(router, h)
	return Test{
		router: router,
		mock:   mockUser,
	}
}

func TestUserHandler_Register(t *testing.T) {
	m := beforeEach(t)
	u, _ := json.Marshal(dto.User{})
	t.Run("Test_for_required_fields", func(t *testing.T) {
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/register", u, resp)
		assert.Contains(t, resp.Body.String(), "validation failed on field 'Email', condition: required")
		assert.Contains(t, resp.Body.String(), "validation failed on field 'FirstName', condition: required")
		assert.Contains(t, resp.Body.String(), "validation failed on field 'LastName', condition: required")
		assert.Contains(t, resp.Body.String(), "validation failed on field 'Password', condition: required")
	})
	user := dto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jdoe@gmail.com",
		Password:  "password",
	}
	u, _ = json.Marshal(user)

	t.Run("Test_for_User_already_exists", func(t *testing.T) {
		m.mock.EXPECT().CheckUserExists(user.Email).Return(user, true)
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/register", u, resp)
		assert.Contains(t, resp.Body.String(), "Email already exists")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("Test_for_failed_registration", func(t *testing.T) {
		m.mock.EXPECT().CheckUserExists(user.Email).Return(user, false)
		m.mock.EXPECT().Create(&user).Return(user, errors.New("an error occurred"))
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/register", u, resp)
		assert.Contains(t, resp.Body.String(), "oops!!!, an error occurred")
	})

	t.Run("Test_for_successful_registration", func(t *testing.T) {
		m.mock.EXPECT().CheckUserExists(user.Email).Return(user, false)
		m.mock.EXPECT().Create(&user).Return(user, nil)
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/register", u, resp)
		assert.Contains(t, resp.Body.String(), "Registration Successful")
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestUSerHandler_Login(t *testing.T) {
	m := beforeEach(t)
	request := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	u, _ := json.Marshal(&request)
	t.Run("Test_for_required_fields", func(t *testing.T) {
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/login", u, resp)
		assert.Contains(t, resp.Body.String(), "validation failed on field 'Email', condition: required")
		assert.Contains(t, resp.Body.String(), "validation failed on field 'Password', condition: required")
	})
	request = struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{
		Email:    "jdoe@gmail.com",
		Password: "password",
	}
	u, _ = json.Marshal(&request)
	t.Run("Test_for_user_not_exists", func(t *testing.T) {
		m.mock.EXPECT().CheckUserExists(request.Email).Return(dto.User{}, false)
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/login", u, resp)
		assert.Contains(t, resp.Body.String(), "Email does not exists")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	user := dto.User{
		Email:        "jdoe@gmail.com",
		HashPassword: "string(hashP)",
	}
	t.Run("Test_for_incorrect_Password", func(t *testing.T) {

		m.mock.EXPECT().CheckUserExists(request.Email).Return(user, true)
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/login", u, resp)
		assert.Contains(t, resp.Body.String(), "Incorrect Password")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	hashP, _ := helpers.GenerateHashPassword(request.Password)
	user = dto.User{
		Email:        "jdoe@gmail.com",
		HashPassword: string(hashP),
	}
	t.Run("Test_for_incorrect_Password", func(t *testing.T) {

		m.mock.EXPECT().CheckUserExists(request.Email).Return(user, true)
		resp := httptest.NewRecorder()
		sendRequest(t, m.router, http.MethodPost, "/api/v1/login", u, resp)
		assert.Contains(t, resp.Body.String(), "Login Successful")
		assert.Contains(t, resp.Body.String(), "token")
		assert.Equal(t, http.StatusOK, resp.Code)
	})

}

func sendRequest(t *testing.T, router *gin.Engine, method string, url string, body []byte,
	resp *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	router.ServeHTTP(resp, req)
}
