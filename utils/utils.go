package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"time"
)

func MkDirIfNotExist(dirPath string) {
	err := os.Mkdir(dirPath, 0755)

	if err != nil && !errors.Is(err, fs.ErrExist) {
		fmt.Println(err)
		os.Exit(1)
	}
}

func WriteBytesDebug(filename string, data []byte) {
	MkDirIfNotExist("debug")

	debugFilepath := path.Join("debug", filename)

	err := os.WriteFile(debugFilepath, data, 0644)

	if err != nil {
		fmt.Printf("cannot write %s: %s", filename, err)

		os.Exit(1)
	}
}

func GetOrdinal(number int) string {
	switch {
	case number%10 == 1 && number%100 != 11:
		{
			return fmt.Sprintf("%dst", number)
		}
	case number%10 == 2 && number%100 != 12:
		{
			return fmt.Sprintf("%dnd", number)
		}
	case number%10 == 3 && number%100 != 13:
		{
			return fmt.Sprintf("%drd", number)
		}
	default:
		{
			return fmt.Sprintf("%dth", number)
		}
	}
}

func FormatDate(time time.Time) string {
	localTime := time.Local()
	monthAndYear := localTime.Format("Jan 2006")
	day := GetOrdinal(localTime.Day())

	return fmt.Sprintf("%s %s", day, monthAndYear)
}
