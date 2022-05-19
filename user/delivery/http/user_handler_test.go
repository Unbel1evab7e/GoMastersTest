package http_test

import (
	"GoMastersTest/models/DTOs"
	"GoMastersTest/models/entity"
	"GoMastersTest/user"
	httpHandler "GoMastersTest/user/delivery/http"
	"GoMastersTest/user/repository"
	"GoMastersTest/user/usecase"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func CustomGenerator() {
	_ = faker.AddProvider("UUID", func(v reflect.Value) (interface{}, error) {
		s, _ := uuid.NewUUID()
		return s, nil
	})
}
func InitRepoAndUseCase() (*user.UseCase, error) {
	viper.SetConfigFile(`../../../config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", dbHost, dbPort, dbName, dbUser, dbPass)
	dsn := fmt.Sprintf("%s", connection)
	dbConn, err := sql.Open(`postgres`, dsn)

	if err != nil {
		return nil, err
	}

	timeoutContext := 60 * time.Second
	repo := repository.NewPostgreUserRepository(dbConn)
	useCase := usecase.NewUserUseCase(repo, timeoutContext)

	return &useCase, nil
}

func TestGetAll(t *testing.T) {
	useCase, err := InitRepoAndUseCase()

	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}

	rec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	ctx, _ := gin.CreateTestContext(rec)
	req = req.WithContext(ctx)
	ctx.Request = req
	handler.GetAllUsers(ctx)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetByID(t *testing.T) {
	CustomGenerator()
	var mockUser entity.User
	err := faker.FakeData(&mockUser)

	assert.NoError(t, err)
	id := mockUser.ID
	useCase, err := InitRepoAndUseCase()

	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}

	rec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v", id), strings.NewReader(""))

	assert.NoError(t, err)
	ctx, _ := gin.CreateTestContext(rec)

	ctx.Params = []gin.Param{
		{
			Key:   "id",
			Value: id.String(),
		},
	}

	ctx.Request = req
	handler.GetUserByID(ctx)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInvalidEmailCreate(t *testing.T) {
	CustomGenerator()
	var userMock DTOs.User

	err := faker.FakeData(&userMock)
	assert.NoError(t, err)

	useCase, err := InitRepoAndUseCase()

	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}

	userJson, err := json.Marshal(userMock)

	assert.NoError(t, err)

	data := url.Values{}
	data.Set("User", string(userJson))

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(string(userJson)))

	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = req

	handler.CreateUser(ctx)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestValidEmailCreate(t *testing.T) {
	CustomGenerator()

	useCase, err := InitRepoAndUseCase()

	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}
	var userMock DTOs.User
	err = faker.FakeData(&userMock)
	userMock.Email = "example@example.com"

	assert.NoError(t, err)

	userJson, err := json.Marshal(userMock)

	assert.NoError(t, err)

	data := url.Values{}
	data.Set("User", string(userJson))

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(string(userJson)))

	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = req
	handler.CreateUser(ctx)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)

}
func TestDelete(t *testing.T) {
	useCase, err := InitRepoAndUseCase()

	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}

	rec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	ctx, _ := gin.CreateTestContext(rec)
	req = req.WithContext(ctx)
	ctx.Request = req
	handler.GetAllUsers(ctx)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	var users []entity.User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	assert.NoError(t, err)

	if len(users) == 0 {
		TestValidEmailCreate(t)
		handler.GetAllUsers(ctx)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &users)
		assert.NoError(t, err)
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(users))
	id := users[i].ID

	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v", id), strings.NewReader(""))

	assert.NoError(t, err)
	ctx, _ = gin.CreateTestContext(rec)

	ctx.Params = []gin.Param{
		{
			Key:   "id",
			Value: id.String(),
		},
	}

	ctx.Request = req
	handler.DeleteUser(ctx)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}
func TestUpdate(t *testing.T) {
	useCase, err := InitRepoAndUseCase()
	var userMock DTOs.User
	err = faker.FakeData(&userMock)
	userMock.Email = "example@example.com"
	assert.NoError(t, err)

	handler := httpHandler.UserHandler{
		Usecase: *useCase,
	}

	rec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	ctx, _ := gin.CreateTestContext(rec)
	req = req.WithContext(ctx)
	ctx.Request = req
	handler.GetAllUsers(ctx)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	var users []entity.User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	assert.NoError(t, err)

	if len(users) == 0 {
		TestValidEmailCreate(t)
		handler.GetAllUsers(ctx)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &users)
		assert.NoError(t, err)
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(users))
	id := users[i].ID

	userJson, err := json.Marshal(userMock)
	assert.NoError(t, err)

	data := url.Values{}
	data.Set("User", string(userJson))

	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%v", id), strings.NewReader(string(userJson)))

	assert.NoError(t, err)
	ctx, _ = gin.CreateTestContext(rec)

	ctx.Params = []gin.Param{
		{
			Key:   "id",
			Value: id.String(),
		},
	}
	ctx.Request = req
	handler.UpdateUser(ctx)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}
