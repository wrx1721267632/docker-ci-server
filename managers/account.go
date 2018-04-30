package managers

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/wrxcode/deploy-server/common/components"
	"github.com/wrxcode/deploy-server/models"
	"github.com/pkg/errors"
)

func AccountLogin(name string, password string) (string, error) {
	account, err := models.Account{}.GetByName(name)
	if err != nil {
		return "", fmt.Errorf("get account failure : %s ", err.Error())
	}
	if account == nil {
		return "", errors.New("name is not exist")
	}
	if account.Password != md5Encode(password) {
		return "", errors.New("Password is wrong")
	}

	token, err := components.CreateToken(account.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func AccountRegister(name, password string) (int64, error) {
	account, err := models.Account{}.GetByName(name)
	if err != nil {
		return 0, fmt.Errorf("get account failure : %s ", err.Error())
	}
	if account != nil {
		return 0, errors.New("name is exist")
	}
	account = &models.Account{Name: name, Password: md5Encode(password)}
	insertId, err := models.Account{}.Add(account)
	if err != nil {
		return 0, fmt.Errorf("add account failure : %s ", err.Error())
	}
	return insertId, nil
}

func md5Encode(password string) string {
	w := md5.New()
	io.WriteString(w, password)
	md5str := string(fmt.Sprintf("%x", w.Sum(nil)))
	return md5str
}

func getCreator(accountId int64) string {
	account, err := models.Account{}.GetById(accountId)
	if err != nil {
		panic(err.Error())
	}
	if account == nil {
		return ""
	}
	return account.Name
}
