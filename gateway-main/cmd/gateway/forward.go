package main

import (
	"context"
	"github.com/BitofferHub/gateway/internal/conf"
	"github.com/BitofferHub/pkg/constant"
	"github.com/BitofferHub/pkg/middlewares/discovery"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Forward(c *gin.Context) {
	action, _ := c.Params.Get("action")
	if _, ok := conf.Routes[action]; !ok {
		c.JSON(http.StatusNotFound, "")
		return
	}
	route := conf.Routes[action]
	userID, _ := c.Get("userID")
	hostReverseProxy(c, c.Writer, c.Request, &route, userID.(string))
}

func hostReverseProxy(ctx context.Context, w http.ResponseWriter, req *http.Request, route *conf.Route, userID string) {
	log.Infof("redirect route: %+v\n", route)
	director := func(req *http.Request) {
		req.Header.Set(constant.UserID, userID)
		req.Header.Set(constant.TraceID, "121321313")
		req.Header.Set("Pika-AAA", "abcdefghhhhhhhh")
		req.URL.Scheme = route.Scheme
		dis := discovery.GetServiceDiscovery(route.Host)
		endpoint := dis.GetHttpEndPoint()
		u, err := url.Parse(endpoint)
		if err != nil {
			panic(err)
		}
		req.URL.Host = u.Host
		req.URL.Path = route.Uri
		log.Infof("redirect url: %+v\n", req.URL)
	}

	log.Infof("redirect host: %+v, url: %+v\n", route.Host, req)
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, req)
}

type TargetHost struct {
	Host string
}
