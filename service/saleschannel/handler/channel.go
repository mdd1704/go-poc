package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	spanLog "github.com/opentracing/opentracing-go/log"
	"github.com/palantir/stacktrace"

	"go-poc/external"
	"go-poc/respond"
	"go-poc/service/saleschannel/model"
	"go-poc/service/saleschannel/usecase"
	"go-poc/utils/activity"
	"go-poc/utils/log"
)

type ChannelHandler struct {
	usecase usecase.Channel
}

func NewChannel(
	usecase usecase.Channel,
) ChannelHandler {
	return ChannelHandler{
		usecase: usecase,
	}
}

func (h *ChannelHandler) HandleUpsert(c *gin.Context) {
	ctx := activity.NewContext("channel_upsert")
	trxID, _ := activity.GetTransactionID(ctx)

	span := external.StartSpanFromRequest(external.Tracer, c.Request, "POST /api/channel/upsert")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	var inputs []model.ChannelInput
	if err := c.BindJSON(&inputs); err != nil {
		span.SetTag("error", true)
		span.LogFields(
			spanLog.String("event", err.Error()),
			spanLog.String("type", respond.ErrBadRequest),
		)

		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	outputs, err := h.usecase.Upsert(ctx, inputs)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			spanLog.String("event", stacktrace.RootCause(err).Error()),
			spanLog.String("type", respond.ErrInternal),
		)

		if len(outputs) > 0 {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error channel upsert %v", outputs))
			respond.Invalid(c, trxID, http.StatusInternalServerError, outputs)
			return
		}

		log.WithContext(ctx).Error("error channel upsert", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	span.LogFields(
		spanLog.String("event", "channel upsert success"),
		spanLog.String("type", "Success"),
	)
	respond.Success(c, trxID, http.StatusCreated, nil)
}

func (h *ChannelHandler) HandleUpsertBatchFetching(c *gin.Context) {
	ctx := activity.NewContext("channel_upsert_batch_fetching")
	trxID, _ := activity.GetTransactionID(ctx)

	span := external.StartSpanFromRequest(external.Tracer, c.Request, "POST /api/channel/upsert-batch-fetching")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	var inputs []model.ChannelInput
	if err := c.BindJSON(&inputs); err != nil {
		span.SetTag("error", true)
		span.LogFields(
			spanLog.String("event", err.Error()),
			spanLog.String("type", respond.ErrBadRequest),
		)

		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	outputs, err := h.usecase.UpsertBatchFetching(ctx, inputs)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			spanLog.String("event", stacktrace.RootCause(err).Error()),
			spanLog.String("type", respond.ErrInternal),
		)

		if len(outputs) > 0 {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error channel upsert %v", outputs))
			respond.Invalid(c, trxID, http.StatusInternalServerError, outputs)
			return
		}

		log.WithContext(ctx).Error("error channel upsert", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	span.LogFields(
		spanLog.String("event", "channel upsert batch fetching success"),
		spanLog.String("type", "Success"),
	)
	respond.Success(c, trxID, http.StatusCreated, nil)
}

func (h *ChannelHandler) HandleUpsertWithTransaction(c *gin.Context) {
	ctx := activity.NewContext("channel_upsert_with_transaction")
	trxID, _ := activity.GetTransactionID(ctx)

	var inputs []model.ChannelInput
	if err := c.BindJSON(&inputs); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	outputs, err := h.usecase.UpsertWithTransaction(ctx, inputs)
	if err != nil {
		if len(outputs) > 0 {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error channel upsert %v", outputs))
			respond.Invalid(c, trxID, http.StatusInternalServerError, outputs)
			return
		}

		log.WithContext(ctx).Error("error channel upsert", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusCreated, nil)
}

func (h *ChannelHandler) HandleUpsertWithLock(c *gin.Context) {
	ctx := activity.NewContext("channel_upsert_with_lock")
	trxID, _ := activity.GetTransactionID(ctx)

	var inputs []model.ChannelInput
	if err := c.BindJSON(&inputs); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	outputs, err := h.usecase.UpsertWithLock(ctx, inputs)
	if err != nil {
		if len(outputs) > 0 {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error channel upsert %v", outputs))
			respond.Invalid(c, trxID, http.StatusInternalServerError, outputs)
			return
		}

		log.WithContext(ctx).Error("error channel upsert", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusCreated, nil)
}

func (h *ChannelHandler) HandleAllByFilter(c *gin.Context) {
	ctx := activity.NewContext("channel_all_by_filter")
	trxID, _ := activity.GetTransactionID(ctx)
	filter := model.ChannelFilter{}
	if err := c.BindJSON(&filter); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	ctx = activity.WithPayload(ctx, filter)
	items, err := h.usecase.FindByFilter(filter)
	if err != nil {
		log.WithContext(ctx).Error("error channel all by filter", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	if len(items) == 0 {
		respond.Error(c, trxID, http.StatusNotFound, respond.ErrNotFound, "channel not found")
		return
	}

	respond.Success(c, trxID, http.StatusOK, items)
}

func (h *ChannelHandler) HandlePagination(c *gin.Context) {
	ctx := activity.NewContext("channel_pagination")
	trxID, _ := activity.GetTransactionID(ctx)

	page := 1
	if number, err := strconv.Atoi(c.Query("page")); err == nil {
		page = number
	}

	limit := 25
	if number, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = number
	}

	filter := model.ChannelFilter{}
	if err := c.BindJSON(&filter); err != nil {
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrBadRequest, err.Error())
		return
	}

	data, err := h.usecase.FindPage(filter, int64(page), int64(limit))
	if err != nil {
		log.WithContext(ctx).Error("error channel pagination", err)
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, data)
}

func (h *ChannelHandler) HandleFindByID(c *gin.Context) {
	ctx := activity.NewContext("channel_find_by_id")
	trxID, _ := activity.GetTransactionID(ctx)

	uri := model.ChannelURI{}
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
		log.WithContext(ctx).Error("error channel find by id", err)
		respond.Error(c, trxID, http.StatusBadRequest, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, data)
}

func (h *ChannelHandler) HandleDelete(c *gin.Context) {
	ctx := activity.NewContext("channel_delete")
	trxID, _ := activity.GetTransactionID(ctx)
	filter := model.ChannelFilter{}
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
		log.WithContext(ctx).Error("error channel delete", err)
		respond.Error(c, trxID, http.StatusInternalServerError, respond.ErrInternal, stacktrace.RootCause(err).Error())
		return
	}

	respond.Success(c, trxID, http.StatusOK, nil)
}
