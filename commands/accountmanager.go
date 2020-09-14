package commands

import (
	"fmt"
)

type accountManager struct {
	user string
	pass string

	Fs *FS
}

func newAccountManager() *accountManager {
	return &accountManager{}
}

func (acc *accountManager) withUser(user string) {
	acc.user = user
	logMsg(fmt.Sprintf("ACC MANAGER SET USER %v", user))
}

func (acc *accountManager) withPass(pass string) {
	acc.pass = pass

	logMsg(fmt.Sprintf("ACC MANAGER SET PASS %v", "******"))
}

func (acc *accountManager) validatePassword(pass string) bool {
	if acc.user == pass {
		acc.withPass(pass)

		return true
	}

	logMsg("PASS CMD user:password validation error")

	return false
}
