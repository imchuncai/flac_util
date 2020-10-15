package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imchuncai/flac"
)

// func main() {
// 	var steam, err = analyze("a.flac")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(steam.VorbisComments)
// 	f, err := os.Create("temp.flac")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	steam.VorbisComments.UserCommentList["ARTIST"] = "imchuncai"

// 	err = steam.Repack(f)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func main() {
// 	fmt.Println("AA")
// 	var steam, err = flac.Analyze("temp.flac")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(steam.VorbisComments)
// }

func main() {
	var f, err = os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	WriteLog(f, "Start")

	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(info.Name(), ".flac") {
			return nil
		}
		var artist, title, ok = ParseFileName(strings.TrimSuffix(info.Name(), ".flac"))
		if !ok {
			WriteLog(f, fmt.Sprintf("The name of file %s is incorrectly formatted", info.Name()))
			return nil
		}
		steam, err := flac.Analyze(path)
		if err != nil {
			WriteLog(f, fmt.Errorf("Analyze %s failed: %w", path, err).Error())
			return nil
		}
		if steam.VorbisComments.UserCommentList["ARTIST"] == artist &&
			steam.VorbisComments.UserCommentList["TITLE"] == title {
			return nil
		}
		steam.VorbisComments.UserCommentList["ARTIST"] = artist
		steam.VorbisComments.UserCommentList["TITLE"] = title
		return steam.RepackFile(path)
	})
	if err != nil {
		fmt.Println(err)
		WriteLog(f, err.Error())
	}
	WriteLog(f, "Finish")
}

func WriteLog(f *os.File, msg string) {
	var _, err = f.Write([]byte(time.Now().Format("2006-01-02 15:04:05") + ": " + msg))
	if err != nil {
		panic(err)
	}
}

func ParseFileName(fileName string) (artist, title string, ok bool) {
	var pieces = strings.Split(fileName, " - ")
	if len(pieces) != 2 {
		return "", "", false
	}
	return pieces[0], pieces[1], true
}
