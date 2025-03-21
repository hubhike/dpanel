package controller

import (
	"errors"
	"github.com/donknap/dpanel/app/common/logic"
	"github.com/donknap/dpanel/common/accessor"
	"github.com/donknap/dpanel/common/dao"
	"github.com/donknap/dpanel/common/entity"
	"github.com/donknap/dpanel/common/service/docker"
	"github.com/donknap/dpanel/common/service/family"
	"github.com/donknap/dpanel/common/service/notice"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/controller"
	"gorm.io/datatypes"
	"gorm.io/gen"
	"time"
)

type User struct {
	controller.Abstract
}

func (self User) Login(http *gin.Context) {
	type ParamsValidate struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		AutoLogin bool   `json:"autoLogin"`
		Code      string `json:"code"`
	}
	params := ParamsValidate{}
	if !self.Validate(http, &params) {
		return
	}

	if new(family.Provider).Check(family.FeatureTwoFa) {
		twoFa := accessor.TwoFa{}
		exists := logic.Setting{}.GetByKey(logic.SettingGroupSetting, logic.SettingGroupSettingTwoFa, &twoFa)
		if exists && twoFa.Enable {
			if params.Code == "" {
				self.JsonResponseWithError(http, errors.New("请输入双因素验证码"), 500)
				return
			}
			if !totp.Validate(params.Code, twoFa.Secret) {
				self.JsonResponseWithError(http, errors.New("验证码错误"), 500)
				return
			}
		}
	}
	var code string
	var err error

	if new(logic.User).CheckLock(params.Username) {
		self.JsonResponseWithError(http, notice.Message{}.New(".userFailedLock", "time", "15"), 500)
		return
	}

	defer func() {
		logic.User{}.Lock(params.Username, code == "")
	}()

	var expireAddTime time.Duration
	if params.AutoLogin {
		expireAddTime = time.Hour * 24 * 30
	} else {
		expireAddTime = time.Hour * 24
	}

	currentUser, err := dao.Setting.Where(dao.Setting.GroupName.Eq(logic.SettingGroupUser)).
		Where(gen.Cond(datatypes.JSONQuery("value").Equals(params.Username, "username"))...).First()
	if err != nil {
		self.JsonResponseWithError(http, notice.Message{}.New(".usernameOrPasswordError"), 500)
		return
	}

	password := logic.User{}.GetMd5Password(params.Password, params.Username)
	if params.Username == currentUser.Value.Username && currentUser.Value.Password == password {
		jwtSecret := logic.User{}.GetJwtSecret()
		jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, logic.UserInfo{
			UserId:       currentUser.ID,
			Username:     currentUser.Value.Username,
			RoleIdentity: currentUser.Name,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireAddTime)),
			},
		})
		code, err = jwtClaims.SignedString(jwtSecret)
		if err != nil {
			self.JsonResponseWithError(http, err, 500)
			return
		}
		self.JsonResponseWithoutError(http, gin.H{
			"accessToken": code,
		})
		return
	} else {
		self.JsonResponseWithError(http, notice.Message{}.New(".usernameOrPasswordError"), 500)
		return
	}
}

func (self User) GetUserInfo(http *gin.Context) {
	result := gin.H{}

	data, exists := http.Get("userInfo")
	if !exists {
		self.JsonResponseWithError(http, errors.New("请先登录"), 401)
		http.AbortWithStatus(401)
		return
	}
	result["user"] = data.(logic.UserInfo)

	feature := []string{
		family.FeatureComposeStore,
	}
	if facade.GetConfig().GetString("app.env") != "lite" && docker.Sdk.Name == docker.DefaultClientName {
		feature = append(feature, family.FeatureContainerDomain)
	}
	result["feature"] = append(feature, family.Provider{}.Feature()...)

	themeConfig := accessor.ThemeConfig{}
	logic.Setting{}.GetByKey(logic.SettingGroupSetting, logic.SettingGroupSettingThemeConfig, &themeConfig)
	result["theme"] = themeConfig

	self.JsonResponseWithoutError(http, result)
	return
}

func (self User) LoginInfo(http *gin.Context) {
	result := gin.H{
		"showRegister":  false,
		"showBuildName": true,
		"family":        facade.GetConfig().GetString("app.env"),
		"feature":       family.Provider{}.Feature(),
		"appName":       facade.GetConfig().GetString("app.name"),
	}
	_, err := logic.Setting{}.GetDPanelInfo()
	if err == nil {
		result["showBuildName"] = false
	}
	_, err = logic.Setting{}.GetValue(logic.SettingGroupUser, logic.SettingGroupUserFounder)
	if err != nil {
		result["showRegister"] = true
	}
	self.JsonResponseWithoutError(http, result)
	return
}

func (self User) SaveThemeConfig(http *gin.Context) {
	params := accessor.ThemeConfig{}
	if !self.Validate(http, &params) {
		return
	}
	err := logic.Setting{}.Save(&entity.Setting{
		GroupName: logic.SettingGroupSetting,
		Name:      logic.SettingGroupSettingThemeConfig,
		Value: &accessor.SettingValueOption{
			ThemeConfig: &params,
		},
	})
	if err != nil {
		self.JsonResponseWithError(http, err, 500)
		return
	}
	self.JsonSuccessResponse(http)
	return
}

func (self User) CreateFounder(http *gin.Context) {
	type ParamsValidate struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}
	params := ParamsValidate{}
	if !self.Validate(http, &params) {
		return
	}
	founder, _ := logic.Setting{}.GetValue(logic.SettingGroupUser, logic.SettingGroupUserFounder)
	if founder != nil {
		self.JsonResponseWithError(http, notice.Message{}.New(".userFounderExists"), 500)
		return
	}
	if params.Password != params.ConfirmPassword {
		self.JsonResponseWithError(http, errors.New(".userPasswordConfirmFailed"), 500)
		return
	}
	err := dao.Setting.Create(&entity.Setting{
		GroupName: logic.SettingGroupUser,
		Name:      logic.SettingGroupUserFounder,
		Value: &accessor.SettingValueOption{
			Password: logic.User{}.GetMd5Password(params.Password, params.Username),
			Username: params.Username,
		},
	})
	if err != nil {
		self.JsonResponseWithError(http, err, 500)
		return
	}
	self.JsonSuccessResponse(http)
	return
}
