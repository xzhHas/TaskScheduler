package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/BitofferHub/user/internal/biz"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}


func (r *userRepo) Save(ctx context.Context, data *biz.Data, g *biz.User) (*biz.User, error) {
	err := data.GetDB().Debug().Create(g).Error
	return g, err
}

func (r *userRepo) Update(ctx context.Context, data *biz.Data, g *biz.User) (*biz.User, error) {
	//err := db.Debug().Update(g).Error
	//return g, err
	return nil, nil
}

// FindByIDWithCache
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description:
//	@Receiver r
//	@param ctx
//	@param data
//	@param userID
//	@return *biz.User
//	@return error
func (r *userRepo) FindByIDWithCache(ctx context.Context, data *biz.Data,
	userID int64) (*biz.User, error) {
	cacheKey := fmt.Sprintf("userinfo:%d", userID)
	var user = new(biz.User)
	rdbUserInfo, exist, err := data.GetCache().Get(ctx, cacheKey)
	if err == nil && exist {
		err = json.Unmarshal([]byte(rdbUserInfo), user)
		if err == nil {
			return user, nil
		}
	}
	user, err = r.FindByID(ctx, data, userID)
	if err != nil {
		return nil, err
	}
	userStr, _ := json.Marshal(user)
	if userStr != nil && len(userStr) != 0 {
		err = data.GetCache().Set(ctx, cacheKey, string(userStr), 10)
		if err != nil {
			log.InfoContextf(ctx, "set user cacheKey err %s", err.Error())
		}
	}
	return user, nil
}

// FindByID
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description:
//	@Receiver r
//	@param ctx
//	@param data
//	@param userID
//	@return *biz.User
//	@return error
func (r *userRepo) FindByID(ctx context.Context, data *biz.Data, userID int64) (*biz.User, error) {
	var user biz.User
	err := data.GetDB().Debug().Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByName
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description:
//	@Receiver r
//	@param ctx
//	@param data
//	@param userName
//	@return *biz.User
//	@return error
func (r *userRepo) FindByName(ctx context.Context, data *biz.Data, userName string) (*biz.User, error) {
	var user biz.User
	err := data.GetDB().Debug().Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListAll
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description:
//	@Receiver r
//	@param ctx
//	@param data
//	@return []*biz.User
//	@return error
func (r *userRepo) ListAll(ctx context.Context, data *biz.Data) ([]*biz.User, error) {
	return nil, nil
}
