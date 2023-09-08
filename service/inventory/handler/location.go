package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"go-poc/respond"
	"go-poc/service/inventory/model"
	"go-poc/service/inventory/usecase"
	"go-poc/utils/activity"
	"go-poc/utils/log"
)

type LocationHandler struct {
	usecase usecase.Location
}

func NewLocation(
	usecase usecase.Location,
) LocationHandler {
	return LocationHandler{
		usecase: usecase,
	}
}

func (h *LocationHandler) HandleUpsert(c *gin.Context) {
	ctx := activity.NewContext("location_upsert")
	trxID, _ := activity.GetTransactionID(ctx)

	var inputs []model.LocationInput
	if err := c.BindJSON(&inputs); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	outputs, err := h.usecase.Upsert(ctx, inputs)
	if err != nil {
		if len(outputs) > 0 {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error location upsert %v", outputs))
			respond.Invalid(c, trxID, http.StatusInternalServerError, outputs)
			return
		}

		log.WithContext(ctx).Error("error location upsert", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusCreated, nil)
}

func (h *LocationHandler) HandleAllByFilter(c *gin.Context) {
	ctx := activity.NewContext("location_all_by_filter")
	trxID, _ := activity.GetTransactionID(ctx)
	filter := model.LocationFilter{}
	if err := c.BindJSON(&filter); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	ctx = activity.WithPayload(ctx, filter)
	items, err := h.usecase.FindByFilter(filter)
	if err != nil {
		log.WithContext(ctx).Error("error location all by filter", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	if len(items) == 0 {
		respond.Error(c, trxID, http.StatusNotFound, respond.ErrNotFound, "location not found")
		return
	}

	respond.Success(c, trxID, http.StatusOK, items)
}

func (h *LocationHandler) HandlePagination(c *gin.Context) {
	ctx := activity.NewContext("location_pagination")
	trxID, _ := activity.GetTransactionID(ctx)

	page := 1
	if number, err := strconv.Atoi(c.Query("page")); err == nil {
		page = number
	}

	limit := 25
	if number, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = number
	}

	filter := model.LocationFilter{}
	if err := c.BindJSON(&filter); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	data, err := h.usecase.FindPage(filter, int64(page), int64(limit))
	if err != nil {
		log.WithContext(ctx).Error("error location pagination", err)
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, data)
}

func (h *LocationHandler) HandleFindByID(c *gin.Context) {
	ctx := activity.NewContext("location_find_by_id")
	trxID, _ := activity.GetTransactionID(ctx)

	uri := model.LocationURI{}
	if err := c.BindUri(&uri); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	data, err := h.usecase.FindByID(id)
	if err != nil {
		log.WithContext(ctx).Error("error location find by id", err)
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, data)
}

func (h *LocationHandler) HandleDelete(c *gin.Context) {
	ctx := activity.NewContext("location_delete")
	trxID, _ := activity.GetTransactionID(ctx)
	filter := model.LocationFilter{}
	if err := c.BindJSON(&filter); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	ctx = activity.WithPayload(ctx, filter)

	if len(filter.IDs) == 0 {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, "ids empty")
		return
	}

	err := h.usecase.Delete(filter)
	if err != nil {
		log.WithContext(ctx).Error("error location delete", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, nil)
}
