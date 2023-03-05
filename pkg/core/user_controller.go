package core

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pergamenum/go-utils-gin/logger"
	m "github.com/pergamenum/go-utils-gin/messages"
	"go.uber.org/zap"
)

// see User.
type userDTO struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name"`
}

// see UserUpdate.
type userUpdateDTO struct {
	ID   string  `json:"id" binding:"required"`
	Name *string `json:"name"`
}

// see UserQuery.
type userQuery url.Values

type UserController struct {
	service UserService
	log     *zap.SugaredLogger
}

func NewUserController(servo UserService) *UserController {
	return &UserController{
		service: servo,
		log:     logger.Get().Named("core.UserController"),
	}
}

// CreateUser - POST: api/v1/user
func (c *UserController) CreateUser(ctx *gin.Context) {

	log := c.log.Named("CreateUser")

	ud := userDTO{}
	err := ctx.ShouldBindJSON(&ud)
	if err != nil {
		log.With("error", err).
			Error("ShouldBindJSON(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user := c.toUser(ud)

	err = c.service.CreateUser(ctx, user)
	switch err {
	case nil:
		ctx.Status(http.StatusCreated)

	case errUserAlreadyExists:
		m.ErrorResponse(ctx, http.StatusConflict, err)

	default:
		log.With("user", user).
			With("error", err).
			Warn("CreateUser(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

// ReadUser - GET: api/v1/user/:id
func (c *UserController) ReadUser(ctx *gin.Context) {

	log := c.log.Named("ReadUser")

	id := ctx.Param("id")

	user, err := c.service.ReadUser(ctx, id)
	switch err {
	case nil:
		ud := c.fromUser(user)
		ctx.JSON(http.StatusOK, ud)

	case errUserNotFound:
		m.ErrorResponse(ctx, http.StatusNotFound, err)

	default:
		log.With("id", id).
			With("error", err).
			Warn("ReadUser(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

// UpdateUser - PATCH: api/v1/user
func (c *UserController) UpdateUser(ctx *gin.Context) {

	log := c.log.Named("UpdateUser")

	uud := userUpdateDTO{}
	err := ctx.ShouldBindJSON(&uud)
	if err != nil {
		log.With("error", err).
			Error("ShouldBindJSON(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	uu := c.toUserUpdate(uud)
	err = c.service.UpdateUser(ctx, uu)
	switch err {
	case nil:
		ctx.Status(http.StatusOK)

	case errUserNotFound:
		m.ErrorResponse(ctx, http.StatusNotFound, err)

	default:
		log.With("uu", uu).
			With("error", err).
			Warn("UpdateUser(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

// DeleteUser - DELETE: api/v1/user/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {

	log := c.log.Named("DeleteUser")

	id := ctx.Param("id")

	err := c.service.DeleteUser(ctx, id)
	switch err {
	case nil:
		ctx.Status(http.StatusOK)

	case errUserNotFound:
		m.ErrorResponse(ctx, http.StatusNotFound, err)

	default:
		log.With("id", id).
			With("error", err).
			Error("...DeleteUser(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

// SearchUsers - GET: api/v1/user
func (c *UserController) SearchUsers(ctx *gin.Context) {

	log := c.log.Named("SearchUsers")

	uq, err := userQuery(ctx.Request.URL.Query()).toUserQuery()
	if err != nil {
		log.With("error", err).
			Error("toUserQuery(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	users, err := c.service.SearchUsers(ctx, uq)
	switch err {
	case nil:
		var uds []userDTO
		for _, user := range users {
			ud := c.fromUser(user)
			uds = append(uds, ud)
		}
		ctx.JSON(http.StatusOK, uds)

	default:
		log.With("error", err).
			Error("SearchUsers(): Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

func (c *UserController) toUser(input userDTO) User {
	return User{
		ID:   input.ID,
		Name: input.Name,
	}
}

func (c *UserController) fromUser(input User) userDTO {
	return userDTO{
		ID:   input.ID,
		Name: input.Name,
	}
}

func (c *UserController) toUserUpdate(input userUpdateDTO) UserUpdate {
	return UserUpdate(input)
}

func (q userQuery) toUserQuery() ([]UserQuery, error) {

	qs, found := q["q"]
	if !found {
		return []UserQuery{}, nil
	}

	var uqs []UserQuery
	for _, q := range qs {
		split := strings.Split(q, ",")
		if len(split) != 3 {
			return nil, errors.New("invalid query: (" + q + ") must be q=(key),(operator),(value)")
		}

		uq := UserQuery{
			Key:      split[0],
			Operator: split[1],
			Value:    split[2],
		}
		uqs = append(uqs, uq)
	}
	return uqs, nil
}
