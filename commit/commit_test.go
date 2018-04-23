package commit

import (
	"testing"
	"flag"
	"fmt"
)

func TestGetCommit(t *testing.T) {
	var gitpath string
	flag.StringVar(&gitpath, "p", "http://github.com/shiyicode/judge", "Git repository")
	ret, err := GetCommit(gitpath)
	fmt.Printf("commit: [%s]\nerror: [%v]", ret, err)
}
