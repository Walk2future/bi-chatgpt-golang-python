package v1

import (
	"github.com/Walk2future/bi-chatgpt-golang-python/common/requests"
	"github.com/Walk2future/bi-chatgpt-golang-python/middleware/jwt"
	"github.com/Walk2future/bi-chatgpt-golang-python/pkg/r"
	"github.com/Walk2future/bi-chatgpt-golang-python/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

var auth = *jwt.AuthMiddleware

// Login godoc
//
//	@Summary	User Login
//	@Produce	json
//	@Tags		UserApi
//	@Param		loginRequest	body	requests.LoginRequest	true	"登录请求参数"
//	@Accept		json
//	@Success	0		{object}	r.Response	"成功"
//	@Failure	40002	{object}	r.Response	"参数错误"
//	@Failure	40003	{object}	r.Response	"系统错误"
//	@Router		/login [post]
func Login(c *gin.Context) {
	auth.LoginHandler(c)
}

// Register godoc
//
//	@Summary	User Register
//	@Produce	json
//	@Tags		UserApi
//	@Param		registerRequest	body	requests.RegisterRequest	true	"注册请求参数"
//	@Accept		json
//	@Success	0		{object}	r.Response	"成功"
//	@Failure	40002	{object}	r.Response	"参数错误"
//	@Failure	40003	{object}	r.Response	"系统错误"
//	@Router		/register [post]
func Register(c *gin.Context) {
	userService := service.UserService{}
	var req requests.RegisterRequest
	validate := validator.New()
	// 使用validator库进行参数校验
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.PARAMS_ERROR.WithMsg("请求参数错误"))
		log.Println(err.Error())
		return
	}
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.SYSTEM_ERROR.WithMsg(err.Error()))
		log.Println(err.Error())
		return
	}
	res, err := userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.SYSTEM_ERROR.WithMsg("注册失败:"+err.Error()))
	} else {
		c.JSON(http.StatusOK, r.OK.WithData(res))
	}
}

// RefreshToken godoc
//
//	@Summary	RefreshToken
//	@Produce	json
//	@Tags		UserApi
//	@Accept		json
//	@Success	0		{object}	r.Response	"成功"
//	@Failure	40005	{object}	r.Response	"认证失败"
//	@Router		/refresh_token [get]
func RefreshToken(c *gin.Context) {
	auth.RefreshHandler(c)
}

// Current godoc
//
//	@Summary	Current
//	@Produce	json
//	@Tags		UserApi
//	@Accept		json
//	@Success	0		{object}	serializers.CurrentUser	"成功"
//	@Failure	40005	{object}	r.Response					"获取当前用户信息失败"
//	@Router		/current [get]
func Current(c *gin.Context) {
	userService := &service.UserService{}
	user := userService.Current(c)
	if user == nil {
		c.JSON(http.StatusBadRequest, r.NO_AUTH.WithMsg("获取当前用户信息失败"))
		c.Abort()
	} else {
		c.JSON(http.StatusOK, r.OK.WithData(user))
	}
}

// Logout godoc
//
//	@Summary	Logout
//	@Produce	json
//	@Tags		UserApi
//	@Accept		json
//	@Success	0		{object}	r.Response	"成功"
//	@Failure	40002	{object}	r.Response	"参数错误"
//	@Failure	40003	{object}	r.Response	"系统错误"
//	@Router		/logout [get]
func Logout(c *gin.Context) {
	auth.LogoutHandler(c)
}
