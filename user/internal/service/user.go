package service

import (
	"context"

	pb "github.com/BitofferHub/user/api/user/v1"
	"github.com/BitofferHub/user/internal/biz"
	"github.com/BitofferHub/user/internal/data"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	dt := biz.NewData(data.GetData().GetDB(), data.GetData().GetCache())
	_, err := s.uc.CreateUser(ctx, dt, &biz.User{
		UserName: req.UserName,
		Pwd:      req.Pwd,
		Sex:      int(req.Sex),
		Age:      int(req.Age),
		Email:    req.Email,
		Contact:  req.Contact,
		Mobile:   req.Mobile,
		IdCard:   req.IdCard,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{Message: "trytest"}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	return &pb.DeleteUserReply{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	dt := biz.NewData(data.GetData().GetDB(), data.GetData().GetCache())
	userInfo, err := s.uc.GetUserInfo(ctx, dt, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserReply{
		Code:    0,
		Message: "success",
		Data: &pb.GetUserReplyData{
			UserName: userInfo.UserName,
			Pwd:      userInfo.Pwd,
			Sex:      int32(userInfo.Sex),
			Age:      int32(userInfo.Age),
			Email:    userInfo.Email,
			Contact:  userInfo.Contact,
			Mobile:   userInfo.Mobile,
			IdCard:   userInfo.IdCard,
		},
	}, nil
}

func (s *UserService) GetUserByName(ctx context.Context, req *pb.GetUserByNameRequest) (*pb.GetUserByNameReply, error) {
	dt := biz.NewData(data.GetData().GetDB(), data.GetData().GetCache())
	userInfo, err := s.uc.GetUserInfoByName(ctx, dt, req.UserName)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByNameReply{
		Code:    0,
		Message: "success",
		Data: &pb.GetUserReplyData{
			UserID:   userInfo.UserID,
			UserName: userInfo.UserName,
			Pwd:      userInfo.Pwd,
			Sex:      int32(userInfo.Sex),
			Age:      int32(userInfo.Age),
			Email:    userInfo.Email,
			Contact:  userInfo.Contact,
			Mobile:   userInfo.Mobile,
			IdCard:   userInfo.IdCard,
		},
	}, nil
}

func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}
