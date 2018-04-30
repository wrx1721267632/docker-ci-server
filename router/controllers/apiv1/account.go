package apiv1

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wrxcode/deploy-server/managers"
	"github.com/wrxcode/deploy-server/router/controllers/base"
)

func RegisterAccount(router *gin.RouterGroup) {
	router.POST("login", httpHandlerLogin)
	router.POST("register", httpHandlerRegister)
}

type AccountParam struct {
	Name     string `form:"name" json:"name"`
	Password string `form:"password" json:"password"`
}

func httpHandlerLogin(c *gin.Context) {
	account := AccountParam{}
	err := c.Bind(&account)
	if err != nil {
		panic(err)
	}
	token, err := managers.AccountLogin(account.Name, account.Password)
	if err != nil {
		c.JSON(http.StatusOK, base.Fail(err.Error()))
		return
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    base64.StdEncoding.EncodeToString([]byte(token)),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerRegister(c *gin.Context) {
	account := AccountParam{}
	err := c.Bind(&account)
	if err != nil {
		panic(err)
	}
	userId, err := managers.AccountRegister(account.Name, account.Password)
	if err != nil {
		c.JSON(http.StatusOK, base.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, base.Success(userId))
}
