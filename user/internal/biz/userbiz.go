package biz

import (
	"context"
	"github.com/BitofferHub/pkg/middlewares/log"
	"time"
)

// User is a User model.
type User struct {
	UserID     string `gorm:"column:id"`
	UserName   string
	Pwd        string
	Sex        int
	Age        int
	Email      string
	Contact    string
	Mobile     string
	IdCard     string
	CreateTime time.Time  `gorm:"column:create_time;default:null"`
	ModifyTime *time.Time `gorm:"column:modify_time;default:null"`
}

// TableName 表名
func (p *User) TableName() string {
	return "t_user_info"
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *Data, *User) (*User, error)
	Update(context.Context, *Data, *User) (*User, error)
	FindByID(context.Context, *Data, int64) (*User, error)
	FindByName(context.Context, *Data, string) (*User, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// CreateUser
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: creates a User, and returns the new User.
//	@Receiver uc
//	@param ctx
//	@param data
//	@param g
//	@return *User
//	@return error
func (uc *UserUsecase) CreateUser(ctx context.Context, data *Data, g *User) (*User, error) {
	log.InfoContextf(ctx, "biz create req username: %v", g.UserName)
	return uc.repo.Save(ctx, data, g)
}

// GetUserInfo
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: get  User, and returns new User.
//	@Receiver uc
//	@param ctx
//	@param data
//	@param userID
//	@return *User
//	@return error
func (uc *UserUsecase) GetUserInfo(ctx context.Context, data *Data, userID int64) (*User, error) {
	return uc.repo.FindByID(ctx, data, userID)
}

// GetUserInfoByName
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: get  User, and returns new User.
//	@Receiver uc
//	@param ctx
//	@param data
//	@param userName
//	@return *User
//	@return error
func (uc *UserUsecase) GetUserInfoByName(ctx context.Context, data *Data, userName string) (*User, error) {
	return uc.repo.FindByName(ctx, data, userName)
}
