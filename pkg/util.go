package pkg

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	Red    = Color("\033[1;31m%s\033[0m")
	Green  = Color("\033[1;32m%s\033[0m")
	Teal   = Color("\033[1;36m%s\033[0m")
	Yellow = Color("\033[1;33m%s\033[0m")
	Purple = Color("\033[1;34m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func ColoredStatus(code int) string {
	switch code {
	case 0:
		return Green("\u2714")
	case 1:
		return Red("\u2717")
	case 15:
		return Yellow("\u2605")
	case 16:
		return White("\u2672")
	case 3012:
		return Red("\u2717")
	default:
		return Teal("\u003f")
	}
	return ""
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func FileContentAsByteArray(file string) []byte {
	fileContent, err := ioutil.ReadFile(file)
	CheckErr(err)
	return fileContent
}

func WriteInFile(file, content string) {
	home, err := homedir.Dir()
	fullPath := home + "/.config/rodai/runs/"
	fileAndPrefix := home + "/.config/rodai/runs/" + file

	err = os.MkdirAll(fullPath, os.ModePerm)
	CheckErr(err)

	f, err := os.OpenFile(fileAndPrefix, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	CheckErr(err)
	defer f.Close()

	f.WriteString(content + "\n")
}

func GenerateFileName(flow string) string {
	now := time.Now().UTC().Unix()
	stringNow := strconv.Itoa(int(now))
	return fmt.Sprintf(flow + "-" + stringNow)
}
