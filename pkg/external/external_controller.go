package external

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	i "github.com/pergamenum/api-gateway/pkg/core/interfaces"
	"github.com/pergamenum/go-utils-gin/logger"
	m "github.com/pergamenum/go-utils-gin/messages"
	"go.uber.org/zap"
)

type Controller[D, M, U any] struct {
	log     *zap.SugaredLogger
	service i.Service[M]
	mapper  i.ControllerMapper[D, M, U]
}

type ControllerConfig[D, M, U any] struct {
	Service i.Service[M]
	Mapper  i.ControllerMapper[D, M, U]
	LogName string
}

func NewController[D, M, U any](conf ControllerConfig[D, M, U]) *Controller[D, M, U] {

	log := logger.Get().Named(conf.LogName)

	c := &Controller[D, M, U]{
		log:     log,
		service: conf.Service,
		mapper:  conf.Mapper,
	}

	return c
}

func (c Controller[D, M, U]) Create(ctx *gin.Context) {

	log := c.log.Named("Create")

	var dto D
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		log.With("error", err).
			Error("ShouldBindJSON: Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	model := c.mapper.ToModel(dto)

	err = c.service.Create(ctx, model)
	switch err {
	case nil:
		ctx.Status(http.StatusCreated)

	case ErrAlreadyExists:
		m.ErrorResponse(ctx, http.StatusConflict, err)

	default:
		log.With("model", model).
			With("error", err).
			Warn("Create: Failed")
		m.ErrorResponse(ctx, http.StatusBadRequest, err)
	}
}

var (
	ErrAlreadyExists = errors.New("resource already exists")
)
