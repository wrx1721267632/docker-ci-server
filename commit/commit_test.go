package commit

import (
	"testing"
	"flag"
	"fmt"
)

func TestGetCommit(t *testing.T) {
	var gitpath string
	flag.StringVar(&gitpath, "p", "github.com/wrxcode/deploy-server", "Git repository")
	ret, _, err := GetCommit(gitpath)
	fmt.Printf("commit: [%s]\nerror: [%v]", ret, err)
}
