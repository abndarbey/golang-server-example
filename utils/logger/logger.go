package logger

import (
	"fmt"
	"log"
	"runtime"
)

// https://pkg.go.dev/github.com/pborman/ansi for reference
const (
	emptySpace     string = " "
	colorReset     string = "\033[0m"
	colorRed       string = "\033[31m"
	colorGreen     string = "\033[32m"
	colorYellow    string = "\033[33m"
	colorBlue      string = "\033[34m"
	colorPurple    string = "\033[35m"
	colorCyan      string = "\033[36m"
	colorWhite     string = "\033[37m"
	boldRedText    string = "\033[1;31m"
	boldGreenText  string = "\033[1;32m"
	boldBlueText   string = "\033[1;34m"
	boldPurpleText string = "\033[1;35m"
	boldCyanText   string = "\033[1;36m"
	boldWhiteText  string = "\033[1;37m"
)

func Info(msg string) {
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	path := fmt.Sprintf("%s:%d", file, line)
	log.Println(
		string(boldBlueText), "Info:",
		string(colorBlue), msg,
		string(boldWhiteText), fn.Name(),
		string(colorReset),
		string(colorWhite), path,
		string(colorReset),
	)
}

func Success(msg string) {
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	path := fmt.Sprintf("%s:%d", file, line)
	log.Println(
		string(boldGreenText), "Success:",
		string(colorGreen), msg,
		string(boldWhiteText), fn.Name(),
		string(colorReset),
		string(colorWhite), path,
		string(colorReset),
	)
}

func Error(err error, msg string) {
	pc, file, line, _ := runtime.Caller(3)
	fn := runtime.FuncForPC(pc)
	path := fmt.Sprintf("%s:%d", file, line)
	log.Println(
		string(boldRedText), "Error:",
		string(colorRed), msg,
		string(colorYellow), err,
		string(boldWhiteText), fn.Name(),
		string(colorReset),
		string(colorWhite), path,
		string(colorReset),
	)
}

func Fatal(msg string) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	path := fmt.Sprintf("%s:%d", file, line)
	log.Fatal(
		string(boldRedText), "Error:",
		emptySpace,
		string(colorRed), msg,
		emptySpace,
		string(boldWhiteText), fn.Name(),
		emptySpace,
		string(colorReset),
		string(colorWhite), path,
		string(colorReset),
	)
}
