package service

import (
	"context"

	"github.com/gin-gonic/gin"

	inventoryHandler "go-poc/service/inventory/handler"
	salesChannelHandler "go-poc/service/saleschannel/handler"
)

func InitRoute(
	ctx context.Context,
	router *gin.Engine,
	channelHandler salesChannelHandler.ChannelHandler,
	locationHandler inventoryHandler.LocationHandler,
) {
	// API group
	api := router.Group("/api")

	api.POST("/channel/upsert", channelHandler.HandleUpsert)
	api.POST("/channel/upsert-batch-fetching", channelHandler.HandleUpsertBatchFetching)
	api.POST("/channel/upsert-with-transaction", channelHandler.HandleUpsertWithTransaction)
	api.POST("/channel/upsert-with-lock", channelHandler.HandleUpsertWithLock)
	api.POST("/channel/filter", channelHandler.HandleAllByFilter)
	api.POST("/channel/pagination", channelHandler.HandlePagination)
	api.DELETE("/channel/delete", channelHandler.HandleDelete)
	api.GET("/channel/:id", channelHandler.HandleFindByID)

	api.POST("/location/upsert", locationHandler.HandleUpsert)
	api.POST("/location/filter", locationHandler.HandleAllByFilter)
	api.POST("/location/pagination", locationHandler.HandlePagination)
	api.DELETE("/location/delete", locationHandler.HandleDelete)
	api.GET("/location/:id", locationHandler.HandleFindByID)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
