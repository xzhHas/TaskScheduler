package service

import (
	"context"
	"encoding/json"

	pb "github.com/BitofferHub/proto_center/api/xtimer/v1"
	"github.com/BitofferHub/xtimer/internal/biz"
	"github.com/BitofferHub/xtimer/internal/constant"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// 创建定时任务
func (s *XTimerService) CreateTimer(ctx context.Context, req *pb.CreateTimerRequest) (*pb.CreateTimerReply, error) {
	param, err := json.Marshal(req.NotifyHTTPParam)
	if err != nil {
		return nil, err
	}
	// 创建定时任
	timer, err := s.timerUC.CreateTimer(ctx, &biz.Timer{
		App:             req.App,
		Name:            req.Name,
		Status:          constant.Unabled.ToInt(),
		Cron:            req.Cron,
		NotifyHTTPParam: string(param),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTimerReply{Code: 0, Message: "ok", Data: &pb.CreateTimerReplyData{
		TimerId: timer.TimerId,
	}}, nil
}

// 激活接口 app + timerId
func (s *XTimerService) EnableTimer(ctx context.Context, req *pb.EnableTimerRequest) (*pb.EnableTimerReply, error) {
	err := s.timerUC.EnableTimer(ctx, req.GetApp(), req.GetTimerId())
	if err != nil {
		return nil, err
	}
	return &pb.EnableTimerReply{Code: 0, Message: "ok"}, nil
}

// 取消激活
func (s *XTimerService) UnableTimer(ctx context.Context, req *pb.EnableTimerRequest) (*pb.EnableTimerReply, error) {
	err := s.timerUC.UnableTimer(ctx, req.GetApp(), req.GetTimerId())
	if err != nil {
		return nil, err
	}
	return &pb.EnableTimerReply{Code: 0, Message: "ok"}, nil
}

// 删除任务
func (s *XTimerService) DleTimer(ctx context.Context, req *pb.EnableTimerRequest) (*pb.EnableTimerReply, error) {
	err := s.timerUC.DelTimer(ctx, req.GetApp(), req.GetTimerId())
	if err != nil {
		return nil, err
	}
	return &pb.EnableTimerReply{Code: 0, Message: "ok"}, nil
}

// 获取所有xtimer
func (s *XTimerService) GetTimers(ctx context.Context) ([]*pb.CreateTimerRequest, error) {
	timers, err := s.timerUC.GetTimers(ctx)
	if err != nil {
		return nil, err
	}

	// 将业务层返回的timers转换为proto格式
	result := make([]*pb.CreateTimerRequest, 0, len(timers))
	for _, timer := range timers {
		// 解析NotifyHTTPParam字符串
		notifyHTTPParam := &pb.NotifyHTTPParam{}
		if timer.NotifyHTTPParam != "" {
			if err := json.Unmarshal([]byte(timer.NotifyHTTPParam), notifyHTTPParam); err != nil {
				// 如果解析失败，记录错误但继续处理
				log.Warnf("解析NotifyHTTPParam失败 %v: %v", timer.TimerId, err)
			}
		}

		// 创建CreateTimerRequest对象
		timerReq := &pb.CreateTimerRequest{
			App:             timer.App,
			Name:            timer.Name,
			Status:          int32(timer.Status),
			Cron:            timer.Cron,
			NotifyHTTPParam: notifyHTTPParam,
		}

		result = append(result, timerReq)
	}

	return result, nil
}
