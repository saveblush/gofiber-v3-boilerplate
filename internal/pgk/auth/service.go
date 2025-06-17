package auth

import (
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/token"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/user"
)

// Service service interface
type Service interface {
	Login(c *cctx.Context, req *RequestLogin) (interface{}, error)
	Logout(c *cctx.Context) error
	LoginBypass(c *cctx.Context, req *RequestLoginBypass) (interface{}, error)
	Bypass(c *cctx.Context, req *RequestBypass) ([]*models.User, error)
	Me(c *cctx.Context) error
	RefreshToken(c *cctx.Context, req *token.RequestRefresh) (*models.Token, error)
}

type service struct {
	config     *config.Configs
	result     *config.ReturnResult
	repository Repository
	cache      cache.Client
	token      token.Service
	user       user.Service
}

func NewService() Service {
	return &service{
		config:     config.CF,
		result:     config.RR,
		repository: NewRepository(),
		cache:      cache.New(),
		token:      token.NewService(),
		user:       user.NewService(),
	}
}

// getLogin get login
func (s *service) getLogin(c *cctx.Context, req *RequestLogLogin) (*models.UserLogin, error) {
	// ยิง log (login สำเร็จ)
	dataLog := &RequestLogLogin{
		UserID:        req.UserID,
		UserLevel:     req.UserLevel,
		EmpID:         req.EmpID,
		ConnectIP:     c.GetClientIP(),
		ConnectDevice: c.GetClientUserAgent(),
		ConnectType:   "login",
		ConnectResult: "success",
	}
	sessionID, err := s.repository.CreateLogLogin(c.GetDatabase(), dataLog)
	if err != nil {
		logger.Log.Error("create log login error: %s", err)
		return nil, err
	}

	//  อัพเดท last login
	dataLast := &RequestLastLogin{
		UserID:        req.UserID,
		ConnectIP:     c.GetClientIP(),
		ConnectDevice: c.GetClientUserAgent(),
	}
	err = s.repository.UpdateLastLogin(c.GetDatabase(), dataLast)
	if err != nil {
		logger.Log.Error("update last login error: %s", err)
		return nil, err
	}

	// ดึงข้อมูลใหม่ user หลังจากที่อัพเดท last login
	fetchUser, err := s.user.Find(c, &user.Request{UserID: req.UserID, UserStatus: []int{1}})
	if err != nil {
		logger.Log.Errorf("find user error: %s", err)
		return nil, err
	}
	if generic.IsEmpty(fetchUser) {
		logger.Log.Errorf("user not found")
		return nil, c.ErrorUnauthorized(s.result.UserNotFound.WithLocale(c.Ctx).Error())
	}

	// create token
	dataToken := &token.Request{
		UserID:    fetchUser.UserID,
		UserLevel: fetchUser.UserLevel,
		EmpID:     fetchUser.EmpID,
		SessionID: sessionID,
	}
	token, err := s.token.Create(c, dataToken)
	if err != nil {
		logger.Log.Errorf("create token error: %s", err)
		return nil, err
	}

	// เก็บ cache session login
	data := &models.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         fetchUser,
	}
	_ = s.createCacheSessionLogin(fetchUser.UserID, sessionID, data)

	res := &models.UserLogin{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         fetchUser,
		SessionID:    sessionID,
	}

	return res, nil
}

// Login login
func (s *service) Login(c *cctx.Context, req *RequestLogin) (interface{}, error) {
	var pass bool

	fetchUser, err := s.user.Find(c, &user.Request{UserID: req.UserID, UserStatus: []int{1}})
	if err != nil {
		logger.Log.Errorf("find user error: %s", err)
		return nil, err
	}

	if !generic.IsEmpty(fetchUser) {
		// check password
		match := comparePassword(fetchUser.Password, req.Password)
		if match {
			pass = true
		}
	}

	if !pass {
		// ยิง log (login ไม่สำเร็จ)
		dataLog := &RequestLogLogin{
			UserID:        req.UserID,
			ConnectIP:     c.GetClientIP(),
			ConnectDevice: c.GetClientUserAgent(),
			ConnectType:   "login",
			ConnectResult: "fail",
		}
		_, err := s.repository.CreateLogLogin(c.GetDatabase(), dataLog)
		if err != nil {
			logger.Log.Error("create log login error: %s", err)
			return nil, err
		}

		return nil, c.ErrorUnauthorized(s.result.InvalidUsername.WithLocale(c.Ctx).Error())
	}

	// login สำเร็จ
	data := &RequestLogLogin{
		UserID:    req.UserID,
		UserLevel: fetchUser.UserLevel,
		EmpID:     fetchUser.EmpID,
	}
	res, err := s.getLogin(c, data)
	if err != nil {
		logger.Log.Error("get login error: %s", err)
		return nil, err
	}

	return res, nil
}

// Logout logout
func (s *service) Logout(c *cctx.Context) error {
	err := s.Me(c)
	if err != nil {
		logger.Log.Errorf("check token not found")
		return c.ErrorUnauthorized(s.result.TokenNotFound.WithLocale(c.Ctx).Error())
	}

	token, err := c.GetClaims()
	if err != nil {
		logger.Log.Errorf("token not found")
		return c.ErrorUnauthorized(s.result.TokenNotFound.WithLocale(c.Ctx).Error())
	}

	// ยิง log logout
	dataLog := &RequestLogLogin{
		UserID:        c.GetUserID(),
		UserLevel:     c.GetUserLevel(),
		EmpID:         c.GetEmpID(),
		ConnectIP:     c.GetClientIP(),
		ConnectDevice: c.GetClientUserAgent(),
		ConnectType:   "logout",
		ConnectResult: "success",
	}
	_, err = s.repository.CreateLogLogin(c.GetDatabase(), dataLog)
	if err != nil {
		logger.Log.Errorf("create log logout error: %s", err)
		return err
	}

	// เคลีย session login
	_ = s.clearCacheSessionLogin(token.Subject, token.SessionID)

	return nil
}

// LoginBypass login bypass
func (s *service) LoginBypass(c *cctx.Context, req *RequestLoginBypass) (interface{}, error) {
	fetchUser, err := s.user.Find(c, &user.Request{UserID: req.UserID, UserStatus: []int{1}})
	if err != nil {
		logger.Log.Errorf("find user error: %s", err)
		return nil, err
	}
	if generic.IsEmpty(fetchUser) {
		logger.Log.Errorf("user not found")
		return nil, c.ErrorUnauthorized(s.result.UserNotFound.WithLocale(c.Ctx).Error())
	}

	data := &RequestLogLogin{
		UserID:    req.UserID,
		UserLevel: fetchUser.UserLevel,
		EmpID:     fetchUser.EmpID,
	}
	res, err := s.getLogin(c, data)
	if err != nil {
		logger.Log.Error("get login error: %s", err)
		return nil, err
	}

	return res, nil
}

// Bypass bypass
func (s *service) Bypass(c *cctx.Context, req *RequestBypass) ([]*models.User, error) {
	// cal api mypass
	//  ส่งค่า req.token เข้าไป จะได้ emp_id กับ user_level กลับมาใช้ดึง user_id
	empID := ""
	userLevel := ""

	res, err := s.user.FindAll(c, &user.Request{EmpID: empID, Userlevel: userLevel, UserStatus: []int{1}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Me me
func (s *service) Me(c *cctx.Context) error {
	token, err := c.GetClaims()
	if err != nil {
		logger.Log.Errorf("token not found")
		return c.ErrorUnauthorized(s.result.TokenNotFound.WithLocale(c.Ctx).Error())
	}

	// เช็ค session login
	err = s.checkSessionLogin(c, token.Subject, token.SessionID, s.getHeaderToken(c))
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken refresh token
func (s *service) RefreshToken(c *cctx.Context, req *token.RequestRefresh) (*models.Token, error) {
	//verify refresh token
	refreshToken, err := s.token.VerifyRefresh(c, req.RefreshToken)
	if err != nil {
		logger.Log.Errorf("verify refreshToken error: %s", err)
		return nil, c.ErrorUnauthorized(s.result.InvalidToken.WithLocale(c.Ctx).Error())
	}

	// renew token
	token, err := s.token.Create(c, &token.Request{
		SessionID: refreshToken.SessionID,
	})
	if err != nil {
		logger.Log.Errorf("create token error: %s", err)
		return nil, s.result.Internal.Unauthorized
	}

	fetchUser, err := s.user.Find(c, &user.Request{UserID: refreshToken.UserID, UserStatus: []int{1}})
	if err != nil {
		logger.Log.Errorf("find user error: %s", err)
		return nil, err
	}
	if generic.IsEmpty(fetchUser) {
		logger.Log.Errorf("user not found")
		return nil, c.ErrorUnauthorized(s.result.UserNotFound.WithLocale(c.Ctx).Error())
	}

	res := &models.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         fetchUser,
	}

	return res, nil
}
