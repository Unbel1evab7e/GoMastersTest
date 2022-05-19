package http

import (
	"GoMastersTest/httputil"
	_ "GoMastersTest/httputil"
	"GoMastersTest/models"
	"GoMastersTest/models/DTOs"
	"GoMastersTest/user"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	Usecase user.UseCase
}

func NewUserHandler(r *gin.Engine, us user.UseCase) {
	handler := &UserHandler{
		Usecase: us,
	}
	v1 := r.Group("/api/v1")
	v1.GET("/users", handler.GetAllUsers)
	v1.GET("/users/:id", handler.GetUserByID)
	v1.POST("/users", handler.CreateUser)
	v1.PUT("/users/:id", handler.UpdateUser)
	v1.DELETE("/users/:id", handler.DeleteUser)

}

// GetAllUsers godoc
// @Summary      Get All Users
// @Description  Return All Users
// @Tags         Users
// @Produce      json
// @Success      200  {object}  []entity.User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Router /users [get]
func (a *UserHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	users, err := a.Usecase.GetAllUsers(ctx)
	if err != nil {
		httputil.NewError(c, getStatusCode(err), err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Get User By id
// @Description  Return User By ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string  true  "User Id"
// @Success      200  {object}  entity.User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Router /users/{id} [get]
func (a *UserHandler) GetUserByID(c *gin.Context) {
	idP := c.Param("id")
	if strings.Count(idP, "")-1 == 0 {
		httputil.NewError(c, http.StatusNotFound, models.ErrNotFound)
		return
	}

	id, err := uuid.Parse(idP)

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, models.ErrBadParamInput)
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := a.Usecase.GetByID(ctx, id)
	if err != nil {
		httputil.NewError(c, getStatusCode(err), err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary CreateUser
// @Description  Return Created User
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        User  body     DTOs.User  true  "Add User"
// @Success      201  {object}  string
// @Failure      422  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Router /users [post]
func (a *UserHandler) CreateUser(c *gin.Context) {
	var user DTOs.User

	jsonDataBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(jsonDataBytes, &user)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if ok, err := isRequestValid(&user); !ok {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	if ok, err := validateEmail(&user); !ok {
		if err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		} else {
			httputil.NewError(c, http.StatusBadRequest, models.ErrEmailValid)
			return
		}
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id, err := a.Usecase.Create(ctx, &user)

	if err != nil {
		httputil.NewError(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusCreated, id.String())
}

// UpdateUser godoc
// @Summary      UpdateUser
// @Description  UpdateUser
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "User ID"
// @Param        User  body     DTOs.User  true  "Update user"
// @Success      200  {object}  entity.User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      422  {object}  httputil.HTTPError
// @Router       /users/{id} [put]
func (a *UserHandler) UpdateUser(c *gin.Context) {
	var user DTOs.User
	idP := c.Param("id")
	if strings.Count(idP, "") == 0 {
		httputil.NewError(c, http.StatusNotFound, models.ErrNotFound)
	}

	id, err := uuid.Parse(idP)

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, models.ErrBadParamInput)
	}

	jsonDataBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(jsonDataBytes, &user)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if ok, err := isRequestValid(&user); !ok {
		httputil.NewError(c, http.StatusBadRequest, err)
	}
	if ok, err := validateEmail(&user); !ok {
		httputil.NewError(c, http.StatusBadRequest, err)
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.Usecase.Update(ctx, id, &user)

	if err != nil {
		httputil.NewError(c, getStatusCode(err), err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteUser godoc
// @Summary 	 DeleteUser
// @Description  DeleteUser
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  string
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Router /users/{id} [delete]
func (a *UserHandler) DeleteUser(c *gin.Context) {
	idP := c.Param("id")
	if strings.Count(idP, "") == 0 {
		httputil.NewError(c, http.StatusNotFound, models.ErrNotFound)
	}

	id, err := uuid.Parse(idP)

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, models.ErrBadParamInput)
	}

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, models.ErrBadParamInput)
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.Usecase.Delete(ctx, id)
	if err != nil {
		httputil.NewError(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusNoContent, "Success")
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
func isRequestValid(m *DTOs.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func validateEmail(userM *DTOs.User) (bool, error) {
	result, err := regexp.Match("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(userM.Email))
	if err != nil {
		return false, err
	}
	return result, nil
}
