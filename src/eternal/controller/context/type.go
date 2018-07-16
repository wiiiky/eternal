package context

import (
	accountModel "eternal/model/account"
	clientModel "eternal/model/client"
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	Token   *accountModel.Token
	Client  *clientModel.Client
	Account *accountModel.Account
}
