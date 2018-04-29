package docker

import (
	"os/exec"
	"bytes"
	"time"
	"fmt"
)

func aaa() {
	cmd := exec.Command("", []string{}...)
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.Stderr = errbuf
	cmd.Stdout = outbuf

	cmd.Start()

	flag := 0

	time.AfterFunc(60*time.Second, func() {
		flag = 3
	})

	go func() {
		err := cmd.Wait()
		if err != nil {
			flag = 2
			fmt.Println("失败", err.Error())
		}

		flag = 1
	}()

	for {
		if flag != 0 {
			break
		}

		out := outbuf.String()
		fmt.Println(out)
	}

	switch flag {
	case 1:
		errout := errbuf.String()
		fmt.Printf(errout)
		if errout == "" {
			fmt.Println("成功")
		} else {
			fmt.Println("失败", errout)
		}

	case 2:
		fmt.Println("失败")
	case 3:
		fmt.Println("超时间")
	}
}
