/**
 * @Description $
 * @Author $
 * @Date $ $
 **/
package utils

import (
	"fmt"
	"os"
	"time"
)

type CommandUtils struct {
}

func (CommandUtils) PauseThenExit(secs int64, val ...any) {
	if len(val) > 0 {
		fmt.Println(val)
	}
	time.Sleep(time.Second * time.Duration(secs))
	os.Exit(-1)
}