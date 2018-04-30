/*
@Time : 18-4-20 下午8:33
@Author : wangruixin
@File : build.go
*/

package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"io"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/wrxcode/deploy-server/models"
)

// 创建docker镜像
func DockerBuild(gitPath string, projectName string, tag string, recordId int64) error {
	gitPath = strings.Replace(gitPath, "http://", "", -1)
	gitPath = strings.Replace(gitPath, "https://", "", -1)

	//获取私有仓库地址
	var repo string
	if g.Conf().Repo.IsHost == 1 {
		repo = g.Conf().Repo.Host
	} else if g.Conf().Repo.IsIp == 1 {
		repo = g.Conf().Repo.Ip
	}
	if repo == "" {
		return errors.Errorf("config file error:[repo]")
	}

	repo = fmt.Sprintf("%s:%s", repo, g.Conf().Repo.Port)

	//拼接镜像名与私有仓库名，方便docker push使用
	repoData := fmt.Sprintf("%s/%s:%s", repo, projectName, tag)

	//执行docker build命令
	cmd := exec.Command("docker", "build", "-t", repoData, gitPath)
	errbuf := new(bytes.Buffer)
	cmd.Stderr = errbuf
	outRead, outWrite := io.Pipe()
	cmd.Stdout = outWrite

	//用来判断docker build是否结束，以及标识结束状态
	flag := 0

	//并发函数，用来接受cmd的结果
	go func() {
		if cmd.Run() == nil {
			flag = 1
		} else {
			flag = 2
		}

		outRead.Close()
	}()

	//定时器任务，当build超过60秒，返回timeout
	time.AfterFunc(600*time.Second, func() {
		cmd.Process.Kill()
		flag = 3
	})

	//持续读入docker build的日志
	scanner := bufio.NewScanner(outRead)
	for scanner.Scan() {
		fmt.Println("===")
		fmt.Printf("%d , rewrite: %s\n", flag, scanner.Text())
		//time.Sleep(1500 * time.Millisecond)
		str, err := rewriteDatabase(recordId, scanner.Text())
		if err != nil {
			log.Fatalf("%s ! construct record id:[%d]; error msg:[%v]", str, recordId, err)
		}
		if flag == 3 {
			fmt.Println("end")
			//fmt.Println(scanner.Text())
			break
		}
		//fmt.Println(scanner.Text())
		//fmt.Println("out-------")
		//
		//data, _ := ioutil.ReadAll(stdout)
		//fmt.Println(string(data))
		//
		//fmt.Println("err-------")
		//data, _ = ioutil.ReadAll(stderr)
		//fmt.Println(string(data))
	}
	fmt.Println("aaaend")
	if flag == 2 {
		errstr := errbuf.String()
		str, err := rewriteDatabase(recordId, errstr)
		if err != nil {
			log.Fatalf("%s ! construct record id:[%d]; error msg:[%v]", str, recordId, err)
		}
		return errors.Errorf("docker build error!!!")
	}

	if flag == 3 {
		str, err := rewriteDatabase(recordId, "\n\ndocker build time out\n")
		if err != nil {
			log.Fatalf("%s ! construct record id:[%d]; error msg:[%v]", str, recordId, err)
		}
		return errors.Errorf("docker build time out!!!")
	}

	return nil
}

// 重写数据库，将build日志重写入数据库
func rewriteDatabase(recordId int64, constructLog string) (string, error) {
	//log.Fatalf("rewrite: %s", constructLog)
	record, err := models.ConstructRecord{}.GetById(recordId)
	if err != nil {
		return "rewriteDatabase: get construct_record error", err
	}

	record.ConstructLog += constructLog
	record.ConstructLog += "\n"
	err = models.ConstructRecord{}.Update(record)
	if err != nil {
		return "rewriteDatabase: update construct_record error", err
	}

	return "", nil
}

// 第一版创建docker镜像函数,已废弃
//func DockerBuild(gitPath string, projectName string, tag string) error {
//gitPath = strings.Replace(gitPath, "http://", "", -1)
//gitPath = strings.Replace(gitPath, "https://", "", -1)
//
//repoData := fmt.Sprintf("%s/%s:%s", g.Conf().DockerHub.Repo, repo, tag)
//
//cmd := exec.Command("docker", "build", "-t", repoData, gitPath)
//stderr := cmd.StderrPipe()
//stdout := cmd.StdoutPipe()
//
//cmd.Start()
//
//flag := 0
//err := errors.New("")
//
//
//go func() {
//	err = cmd.Wait()
//	out := out_buf.String()
//	fmt.Println(out)
//	if err != nil {
//		flag = 2
//		fmt.Println("失败", err.Error())
//	}
//	reader := bufio.NewReader(stdout)
//	reader.
//	flag = 1
//}()
//
//
//
//switch flag {
//case 1:
//	errout := err_buf.String()
//	fmt.Printf(errout)
//	if errout == "" {
//		fmt.Println("成功")
//		return nil
//	} else {
//		fmt.Println("失败", errout)
//		return errors.Errorf("fail")
//	}
//
//case 2:
//	fmt.Println("失败")
//	return err
//case 3:
//	fmt.Println("超时间")
//	return errors.Errorf("time out")
//}
//
//
//return nil
//	err := cmd.Run()
//	if err != nil {
//		return err_out.String(), err
//	}
//	if err_out.String() != "" {
//		return err_out.String(), errors.Errorf("docker build exec error")
//	}
//
//	return out.String(), nil
// }
