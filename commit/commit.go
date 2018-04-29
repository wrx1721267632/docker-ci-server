/*
@Time : 18-4-20 下午5:23 
@Author : wangruixin
@File : commit.go
*/

package commit

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"fmt"
	"strings"
)

// 通过git爬取到commit的信息
func GetCommit(gitPath string) (string, string, error){
	gitPath = strings.Replace(gitPath, "http://", "", -1)
	gitPath = strings.Replace(gitPath, "https://", "", -1)

	if gitPath[len(gitPath) - 1] != '/' {
		gitPath += "/"
	}

	//拼接gitPath 获得存有commit信息的页面
	gitUrl := "https://"
	gitUrl += gitPath
	gitUrl += "commits/master"

	//http GET请求获取对应url中的response对象
	//response, err := http.Post(gitUrl, "application/x-www-form-urlencoded", strings.NewReader("name=xxx"))
	response, err := http.Get(gitUrl)
	//如果访问不成功或者url不存在则会进入改判断
	if err != nil {
		return "GET Error", "", err
	}

	//请求完了关闭回复主体
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "ReadBody Error", "", err
	}

	//commitReg := fmt.Sprintf("(?<=%scommit/).+?(?=\")", gitPath)
	commitReg := fmt.Sprintf("%scommit/(.+?)\"", gitPath)
	reg, err := regexp.Compile(commitReg)
	if err != nil {
		return commitReg, "", err
	}

	retCommit := reg.FindString(string(body))
	if retCommit == "" {
		return "Commit not find", "", fmt.Errorf("regexp get nil: regexp[%s]", commitReg)
	}

	retCommit = strings.Replace(retCommit, fmt.Sprintf("%scommit/", gitPath), "", 1)
	retCommit = strings.Replace(retCommit, "\"", "", 1)

	descReg := "a title=\"(.+?)\""
	reg_desc, err := regexp.Compile(descReg)
	if err != nil {
		return descReg, "", err
	}

	retDesc := reg_desc.FindString(string(body))
	if retDesc == "" {
		return "Commit not find", "", fmt.Errorf("regexp get nil: regexp[%s]", descReg)
	}

	return retCommit, retDesc, nil
}