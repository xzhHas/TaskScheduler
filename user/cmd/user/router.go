package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BitofferHub/pkg/constant"
	"github.com/BitofferHub/pkg/middlewares/log"
	pb "github.com/BitofferHub/user/api/user/v1"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	nethttp "net/http"
	"strconv"
	"time"
)

func GetUserInfo(ctx *gin.Context) {
	userIDStr := ctx.Request.Header.Get(constant.UserID)
	traceID := ctx.Request.Header.Get(constant.TraceID)

	userID, _ := strconv.Atoi(userIDStr)
	//userName := fmt.Sprintf("%s", userID)
	var req = pb.GetUserRequest{
		UserID: int64(userID),
	}
	c := context.WithValue(context.Background(), constant.TraceID, traceID)
	resp, err := us.GetUser(c, &req)
	if err != nil {
		fmt.Println("get user err", err)
	}
	ctx.JSON(nethttp.StatusOK, resp)
}

// InfoLog
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: gin middleware for log request and reply
//	@return gin.HandlerFunc
func InfoLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()
		// ***** 1. get request body ****** //
		traceID := c.Request.Header.Get(constant.TraceID)
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close() //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// ***** 2. set requestID for goroutine ctx ****** //
		//duration := float64(time.Since(beginTime)) / float64(time.Second)
		ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
		log.InfoContextf(ctx, "ReqPath[%s]-Cost[%v]\n", c.Request.URL.Path, time.Since(beginTime))
	}
}
