package interfaces

import (
	"context"
	"fmt"

	"github.com/BitofferHub/pkg/constant"
	"github.com/BitofferHub/pkg/middlewares/log"
	pb "github.com/BitofferHub/proto_center/api/xtimer/v1"
	"github.com/BitofferHub/xtimer/internal/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTimer(c *gin.Context) {
	traceID := c.Request.Header.Get(constant.TraceID)

	var req pb.CreateTimerRequest

	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("CreateTimer err: %v", err)
		response.Fail(c, response.ParamError, nil)
		return
	}

	ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
	// CreaterTimer 创建定时任务的服务层调用接口
	resp, err := h.xTimerService.CreateTimer(ctx, &req)
	if err != nil {
		fmt.Println("Create timer err: %v", err)
		response.Fail(c, response.ParamError, nil)
		return
	}

	response.Success(c, resp)
}

func (h *Handler) EnableTimer(c *gin.Context) {
	traceID := c.Request.Header.Get(constant.TraceID)
	// timerId, err := strconv.ParseInt(c.Query("timerId"), 10, 64)
	// if err != nil {
	// 	log.Errorf("EnableTimer err: %v", err)
	// 	response.Fail(c, response.ParamError, nil)
	// 	return
	// }
	req := pb.EnableTimerRequest{
		//TimerId: timerId,
		App: c.Query("app"),
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.ParamError, nil)
		return
	}

	ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
	_, err := h.xTimerService.EnableTimer(ctx, &req)
	if err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.EnableTimerError, nil)
		return

	}
	response.Success(c, nil)
	return
}

func (h *Handler) UnableTimer(c *gin.Context) {
	traceID := c.Request.Header.Get(constant.TraceID)
	// timerId, err := strconv.ParseInt(c.Query("timerId"), 10, 64)
	// if err != nil {
	// 	log.Errorf("UnableTimer err: %v", err)
	// 	response.Fail(c, response.ParamError, nil)
	// 	return
	// }
	req := pb.EnableTimerRequest{
		TimerId: int64(1),
		App:     c.Query("app"),
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.ParamError, nil)
		return
	}

	ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
	_, err := h.xTimerService.UnableTimer(ctx, &req)
	if err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.EnableTimerError, nil)
		return

	}
	response.Success(c, nil)
	return
}
func (h *Handler) DelTimer(c *gin.Context) {
	traceID := c.Request.Header.Get(constant.TraceID)
	//timerId, err := strconv.ParseInt(c.Query("timerId"), 10, 64)

	req := pb.EnableTimerRequest{
		TimerId: int64(1),
		App:     c.Query("app"),
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.ParamError, nil)
		return
	}

	ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
	_, err := h.xTimerService.DleTimer(ctx, &req)
	if err != nil {
		log.Errorf("EnableTimer err: %v", err)
		response.Fail(c, response.EnableTimerError, nil)
		return

	}
	response.Success(c, nil)
	return
}

func (h *Handler) GetTimers(c *gin.Context) {
	traceID := c.Request.Header.Get(constant.TraceID)

	ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
	timers, err := h.xTimerService.GetTimers(ctx)
	if err != nil {
		log.Errorf("GetTimers err: %v", err)
		response.Fail(c, response.GetTimersError, nil)
		return
	}

	// 将结果包装成带有状态码的响应
	result := map[string]interface{}{
		"code":    0,
		"message": "ok",
		"data":    timers,
	}

	response.Success(c, result)
}

func (h *Handler) TestCallback(c *gin.Context) {
	log.Info("callback test: %v", c.Request.Body)

	response.Success(c, "ok: callback receives")
}
