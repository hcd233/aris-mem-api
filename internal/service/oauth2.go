package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/jwt"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/resource/database"
	"github.com/hcd233/aris-mem-api/internal/resource/database/dao"
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
	objdao "github.com/hcd233/aris-mem-api/internal/resource/storage/obj_dao"

	"github.com/hcd233/aris-mem-api/internal/oauth2"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Oauth2Service OAuth2服务接口
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type Oauth2Service interface {
	Login(ctx context.Context, req *dto.LoginReq) (rsp *dto.LoginResp, err error)
	Callback(ctx context.Context, req *dto.CallbackReq) (rsp *dto.CallbackResp, err error)
}

// oauth2Service OAuth2服务基础实现
type oauth2Service struct {
	provider           oauth2.Provider
	userDAO            *dao.UserDAO
	imageObjDAO        objdao.ObjDAO
	thumbnailObjDAO    objdao.ObjDAO
	accessTokenSigner  jwt.TokenSigner
	refreshTokenSigner jwt.TokenSigner
}

// NewGithubOauth2Service 创建Github OAuth2服务
func NewGithubOauth2Service() Oauth2Service {
	return &oauth2Service{
		provider: oauth2.NewGithubProvider(),
		userDAO:  dao.GetUserDAO(),
		// imageObjDAO:        objdao.GetImageObjDAO(),
		// thumbnailObjDAO:    objdao.GetThumbnailObjDAO(),
		accessTokenSigner:  jwt.GetAccessTokenSigner(),
		refreshTokenSigner: jwt.GetRefreshTokenSigner(),
	}
}

// NewGoogleOauth2Service 创建Google OAuth2服务
func NewGoogleOauth2Service() Oauth2Service {
	return &oauth2Service{
		provider: oauth2.NewGoogleProvider(),
		userDAO:  dao.GetUserDAO(),
		// imageObjDAO:        objdao.GetImageObjDAO(),
		// thumbnailObjDAO:    objdao.GetThumbnailObjDAO(),
		accessTokenSigner:  jwt.GetAccessTokenSigner(),
		refreshTokenSigner: jwt.GetRefreshTokenSigner(),
	}
}

// Login 登录
//
//	receiver s *oauth2Service
//	param ctx context.Context
//	param req *dto.LoginRequest
//	return rsp *dto.LoginResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (s *oauth2Service) Login(ctx context.Context, req *dto.LoginReq) (rsp *dto.LoginResp, err error) {
	rsp = &dto.LoginResp{}

	logger := logger.WithCtx(ctx)

	url := s.provider.GetAuthURL()
	rsp.RedirectURL = url

	logger.Info("[Oauth2Service] login", zap.String("provider", req.Provider), zap.String("redirectURL", url))

	return rsp, nil
}

// Callback 回调
//
//	receiver s *oauth2Service
//	param ctx context.Context
//	param req *dto.CallbackRequest
//	return rsp *dto.CallbackResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (s *oauth2Service) Callback(ctx context.Context, req *dto.CallbackReq) (rsp *dto.CallbackResp, err error) {
	rsp = &dto.CallbackResp{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req.State != config.Oauth2StateString {
		logger.Error("[Oauth2Service] invalid state",
			zap.String("provider", req.Provider),
			zap.String("state", req.State),
			zap.String("expectedState", config.Oauth2StateString))
		return nil, constant.ErrUnauthorized
	}

	logger.Info("[Oauth2Service] exchanging token",
		zap.String("provider", req.Provider),
		zap.String("code", req.Code),
		zap.String("state", req.State))

	token, err := s.provider.ExchangeToken(ctx, req.Code)
	if err != nil {
		logger.Error("[Oauth2Service] failed to exchange token",
			zap.String("provider", req.Provider),
			zap.String("code", req.Code),
			zap.Error(err))
		return nil, constant.ErrUnauthorized
	}

	logger.Info("[Oauth2Service] token exchange successful",
		zap.String("provider", req.Provider),
		zap.String("tokenType", token.TokenType),
		zap.Bool("valid", token.Valid()))

	userInfo, err := s.provider.GetUserInfo(ctx, token)
	if err != nil {
		logger.Error("[Oauth2Service] failed to get user info",
			zap.String("provider", req.Provider),
			zap.Error(err))
		return nil, constant.ErrInternalError
	}

	thirdPartyID := userInfo.GetID()
	userName, email, avatar := userInfo.GetName(), userInfo.GetEmail(), userInfo.GetAvatar()

	user, err := s.userDAO.GetByEmail(db, email, []string{"id", "name", "avatar"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[Oauth2Service] failed to get user by email",
			zap.String("provider", req.Provider),
			zap.String("email", email),
			zap.Error(err))
		return nil, constant.ErrInternalError
	}

	if user.ID != 0 {
		// 更新已存在用户的登录时间
		if err := s.userDAO.Update(db, user, map[string]interface{}{
			"last_login": time.Now().UTC(),
		}); err != nil {
			logger.Error("[Oauth2Service] failed to update user login time",
				zap.String("provider", req.Provider),
				zap.Error(err))
			return nil, constant.ErrInternalError
		}
	} else {
		// 创建新用户
		if validateErr := util.ValidateUserName(userName); validateErr != nil {
			userName = "ArisUser" + strconv.FormatInt(time.Now().UTC().Unix(), 10)
		}
		user = &model.User{
			Name:       userName,
			Email:      email,
			Avatar:     avatar,
			Permission: enum.PermissionReader,
			LastLogin:  time.Now().UTC(),
		}

		if err := s.userDAO.Create(db, user); err != nil {
			logger.Error("[Oauth2Service] failed to create user",
				zap.String("provider", req.Provider),
				zap.String("userName", userName),
				zap.Error(err))
			return nil, constant.ErrInternalError
		}

		// _, err = s.imageObjDAO.CreateDir(ctx, user.ID)
		// if err != nil {
		// 	logger.Error("[Oauth2Service] failed to create image dir",
		// 		zap.String("provider", req.Provider),
		// 		zap.Error(err))
		// 	return nil, protocol.ErrInternalError
		// }
		// logger.Info("[Oauth2Service] image dir created", zap.String("provider", req.Provider))

		// _, err = s.thumbnailObjDAO.CreateDir(ctx, user.ID)
		// if err != nil {
		// 	logger.Error("[Oauth2Service] failed to create thumbnail dir",
		// 		zap.String("provider", req.Provider),
		// 		zap.Error(err))
		// 	return nil, protocol.ErrInternalError
		// }
		// logger.Info("[Oauth2Service] thumbnail dir created", zap.String("provider", req.Provider))
	}

	// 更新第三方平台绑定ID
	bindField := s.provider.GetBindField()
	updateData := map[string]interface{}{
		bindField: thirdPartyID,
	}

	if err := s.userDAO.Update(db, user, updateData); err != nil {
		logger.Error("[Oauth2Service] failed to update third party bind id",
			zap.String("provider", req.Provider),
			zap.String("bindField", bindField),
			zap.String("thirdPartyID", thirdPartyID),
			zap.Error(err))
		return nil, constant.ErrInternalError
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode access token",
			zap.String("provider", req.Provider),
			zap.Error(err))
		return nil, constant.ErrInternalError
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode refresh token",
			zap.String("provider", req.Provider),
			zap.Error(err))
		return nil, constant.ErrInternalError
	}

	logger.Info("[Oauth2Service] callback success",
		zap.String("provider", req.Provider),
		zap.Uint("userID", user.ID))

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
