package commands

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
