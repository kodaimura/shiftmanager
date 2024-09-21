package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/errs"
	"shiftmanager/internal/service"
)

type AccountController struct {
	accountService service.AccountService
}


func NewAccountController() *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(),
	}
}


//GET /signup
func (ctr *AccountController) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

//GET /login
func (ctr *AccountController) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}


//GET /logout
func (ctr *AccountController) Logout(c *gin.Context) {
	jwt.RemoveTokenFromCookie(c)
	c.Redirect(303, "/login")
}


//POST /api/signup
func (ctr *AccountController) ApiSignup(c *gin.Context) {
	m := map[string]string{}
	c.BindJSON(&m)
	name := m["account_name"]
	pass := m["account_password"]

	err := ctr.accountService.Signup(name, pass)
	if err != nil {
		if _, ok := err.(errs.UniqueConstraintError); ok {
			c.JSON(409, gin.H{"error": "ユーザ名が既に使われています。"})
		} else {
			c.JSON(500, gin.H{"error": "登録に失敗しました。"})
		}
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}


//POST /api/login
func (ctr *AccountController) ApiLogin(c *gin.Context) {
	m := map[string]string{}
	c.BindJSON(&m)
	name := m["account_name"]
	pass := m["account_password"]

	account, err := ctr.accountService.Login(name, pass)
	if err != nil {
		c.JSON(401, gin.H{"error": "ユーザ名またはパスワードが異なります。"})
		c.Abort()
		return
	}

	pl, err := ctr.accountService.GenerateJwtPayload(account.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": "ログインに失敗しました。"})
		c.Abort()
		return
	}

	jwt.SetTokenToCookie(c, pl)
	c.JSON(200, gin.H{})
}


//GET /api/account/profile
func (ctr *AccountController) ApiGetProfile(c *gin.Context) {
	pl := jwt.GetPayload(c)
	account, err := ctr.accountService.GetProfile(pl.AccountId)

	if err != nil {
		c.JSON(500, gin.H{})
		c.Abort()
		return
	}

	c.JSON(200, account)
}


//PUT /api/account/password
func (ctr *AccountController) ApiPutPassword(c *gin.Context) {
	pl := jwt.GetPayload(c)
	id := pl.AccountId
	name := pl.AccountName

	m := map[string]string{}
	c.BindJSON(&m)
	oldPass := m["old_account_password"]
	pass := m["account_password"]

	_, err := ctr.accountService.Login(name, oldPass)
	if err != nil {
		c.JSON(400, gin.H{"error": "旧パスワードが異なります。"})
		c.Abort()
		return
	}

	if ctr.accountService.UpdatePassword(id, pass) != nil {
		c.JSON(500, gin.H{"error": "変更に失敗しました。"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}


//PUT /api/account/name
func (ctr *AccountController) ApiPutName(c *gin.Context) {
	pl := jwt.GetPayload(c)
	id := pl.AccountId

	m := map[string]string{}
	c.BindJSON(&m)
	name := m["account_name"]

	err := ctr.accountService.UpdateName(id, name)
	if err != nil {
		if _, ok := err.(errs.UniqueConstraintError); ok {
			c.JSON(409, gin.H{"error": "ユーザ名が既に使われています。"})
		} else {
			c.JSON(500, gin.H{"error": "変更に失敗しました。"})
		}
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}


//DELETE /api/account
func (ctr *AccountController) ApiDeleteAccount(c *gin.Context) {
	pl := jwt.GetPayload(c)
	id := pl.AccountId

	if ctr.accountService.DeleteAccount(id) != nil {
		c.JSON(500, gin.H{"error": "削除に失敗しました。"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}