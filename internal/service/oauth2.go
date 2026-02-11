package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/pool"
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/jwt"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/samber/lo"

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
	Callback(ctx context.Context, req *dto.CallbackReq) (rsp *dto.CallbackRsp, err error)
}

// oauth2Service OAuth2服务基础实现
type oauth2Service struct {
	platform           oauth2.Platform
	userDAO            *dao.UserDAO
	audioObjDAO        objdao.ObjDAO
	accessTokenSigner  jwt.TokenSigner
	refreshTokenSigner jwt.TokenSigner
}

// NewGithubOauth2Service 创建Github OAuth2服务
func NewGithubOauth2Service() Oauth2Service {
	return &oauth2Service{
		platform:           oauth2.NewGithubPlatform(),
		userDAO:            dao.GetUserDAO(),
		audioObjDAO:        objdao.GetAudioObjDAO(),
		accessTokenSigner:  jwt.GetAccessTokenSigner(),
		refreshTokenSigner: jwt.GetRefreshTokenSigner(),
	}
}

// NewGoogleOauth2Service 创建Google OAuth2服务
func NewGoogleOauth2Service() Oauth2Service {
	return &oauth2Service{
		platform:           oauth2.NewGooglePlatform(),
		userDAO:            dao.GetUserDAO(),
		audioObjDAO:        objdao.GetAudioObjDAO(),
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

	url := s.platform.GetAuthURL()
	rsp.RedirectURL = url

	logger.Info("[Oauth2Service] login", zap.String("platform", req.Platform), zap.String("redirectURL", url))

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
func (s *oauth2Service) Callback(ctx context.Context, req *dto.CallbackReq) (*dto.CallbackRsp, error) {
	rsp := &dto.CallbackRsp{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req.Body.State != config.Oauth2StateString {
		logger.Error("[Oauth2Service] invalid state",
			zap.String("platform", req.Body.Platform),
			zap.String("state", req.Body.State),
			zap.String("expectedState", config.Oauth2StateString))
		rsp.Error = constant.ErrUnauthorized
		return rsp, nil
	}

	logger.Info("[Oauth2Service] exchanging token",
		zap.String("platform", req.Body.Platform),
		zap.String("code", req.Body.Code),
		zap.String("state", req.Body.State))

	token, err := s.platform.ExchangeToken(ctx, req.Body.Code)
	if err != nil {
		logger.Error("[Oauth2Service] failed to exchange token",
			zap.String("platform", req.Body.Platform),
			zap.String("code", req.Body.Code),
			zap.Error(err))
		rsp.Error = constant.ErrUnauthorized
		return rsp, nil
	}

	logger.Info("[Oauth2Service] token exchange successful",
		zap.String("platform", req.Body.Platform),
		zap.String("tokenType", token.TokenType),
		zap.Bool("valid", token.Valid()))

	userInfo, err := s.platform.GetUserInfo(ctx, token)
	if err != nil {
		logger.Error("[Oauth2Service] failed to get user info",
			zap.String("platform", req.Body.Platform),
			zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	thirdPartyID := userInfo.GetID()
	userName, email, avatar := userInfo.GetName(), userInfo.GetEmail(), userInfo.GetAvatar()
	if thirdPartyID == "" || thirdPartyID == "0" {
		logger.Error("[Oauth2Service] invalid third party id", zap.String("platform", req.Body.Platform), zap.String("thirdPartyID", thirdPartyID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	var user *model.User
	switch req.Body.Platform {
	case enum.Oauth2PlatformGithub:
		user, err = s.userDAO.Get(db, &model.User{GithubBindID: thirdPartyID}, []string{"id"})
	case enum.Oauth2PlatformGoogle:
		user, err = s.userDAO.Get(db, &model.User{GoogleBindID: thirdPartyID}, []string{"id"})
	default:
		logger.Error("[Oauth2Service] invalid platform", zap.String("platform", req.Body.Platform))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[Oauth2Service] failed to get user by third party bind id",
			zap.String("platform", req.Body.Platform),
			zap.String("thirdPartyID", thirdPartyID),
			zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if user.ID != 0 {
		// 更新已存在用户的登录时间
		if err := s.userDAO.Update(db, user, map[string]interface{}{
			"last_login": time.Now().UTC(),
		}); err != nil {
			logger.Error("[Oauth2Service] failed to update user login time",
				zap.String("platform", req.Body.Platform),
				zap.Error(err))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
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
			Permission: enum.PermissionPending,
			LastLogin:  time.Now().UTC(),
		}

		switch req.Body.Platform {
		case enum.Oauth2PlatformGithub:
			user.GithubBindID = thirdPartyID
		case enum.Oauth2PlatformGoogle:
			user.GoogleBindID = thirdPartyID
		}

		if err := s.userDAO.Create(db, user); err != nil {
			logger.Error("[Oauth2Service] failed to create user",
				zap.String("platform", req.Body.Platform),
				zap.String("userName", userName),
				zap.Error(err))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		_, err = s.audioObjDAO.CreateDir(ctx, user.ID)
		if err != nil {
			logger.Error("[Oauth2Service] failed to create audio dir",
				zap.String("platform", req.Body.Platform),
				zap.Error(err))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		logger.Info("[Oauth2Service] audio dir created", zap.String("platform", req.Body.Platform))

		// 查询前5个管理员用户
		var adminUsers []*model.User
		adminUsers, _, err := s.userDAO.Paginate(db, &model.User{Permission: enum.PermissionAdmin}, []string{"id", "email", "name"}, &dao.CommonParam{
			PageParam: dao.PageParam{
				Page:     1,
				PageSize: 5,
			},
		})
		if err != nil {
			logger.Warn("[Oauth2Service] failed to get admin users for notification",
				zap.Error(err))
		} else if len(adminUsers) > 0 {

			adminEmails := lo.Map(adminUsers, func(item *model.User, _ int) string {
				return item.Email
			})
			if len(adminEmails) > 0 {
				// 构建邮件内容
				subject := fmt.Sprintf("New User Registration: %s", user.Name)
				htmlBody := fmt.Sprintf(constant.NewUserRegistrationEmailTemplate, user.Name, user.Email, user.Avatar, user.CreatedAt.Format(time.RFC3339), user.ID, config.ServerEndpoint)

				task := &dto.EmailSendTask{
					Ctx:      util.CopyContextValues(ctx),
					Emails:   adminEmails,
					Subject:  subject,
					HTMLBody: htmlBody,
				}

				poolManager := pool.GetPoolManager()
				// 提交邮件通知任务
				if err := poolManager.SubmitEmailSendTask(task); err != nil {
					logger.Warn("[Oauth2Service] failed to submit email notification task",
						zap.Error(err),
						zap.Uint("newUserID", user.ID))
				}
			}
		}
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode access token",
			zap.String("platform", req.Body.Platform),
			zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode refresh token",
			zap.String("platform", req.Body.Platform),
			zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	logger.Info("[Oauth2Service] callback success",
		zap.String("platform", req.Body.Platform),
		zap.Uint("userID", user.ID))

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
