package context

import (
	accountModel "eternal/model/account"
	clientModel "eternal/model/client"
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	Client  *clientModel.Client
	Account *accountModel.Account
}
