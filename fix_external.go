package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func fixExternalRefs(htmlFn string, htmlDir string, dir string) error {
	bytes, err := os.ReadFile(filepath.Join(htmlDir, htmlFn))
	if err != nil {
		return err
	}
	str := string(bytes)

	regex := regexp.MustCompile("\"(?U)(https://(\\S+))\"")
	matches := regex.FindAllStringSubmatch(str, -1)

	for i := range matches {
		url := strings.TrimSuffix(matches[i][1], "\"")
		unsafe := matches[i][2]
		badChars := []string{"#", "$", "%", "!", "&", "'", "{", "\"", "}", ":", "\\", "@", "<", "+", ">", "`", "*", "|", "?", "=", "_", ";", "~"}
		for x := range badChars {
			unsafe = strings.ReplaceAll(unsafe, badChars[x], "-")
		}
		fn := strings.ReplaceAll(unsafe, "/", "_")

		if len(dir) > 248 {
			return errors.New("depth limit reached")
		}
		if len(fn) > 250-len(dir) {
			fn = fn[int(math.Abs(float64(249-len(fn)-len(dir)))):]
		}
		fp := filepath.Join(dir, fn)
		_, er := os.Stat(fp)
		if er != nil {
			if errors.Is(er, os.ErrNotExist) {
				if isBlocked(url) {
					continue
				}

				resp, e := http.Get(url)
				if e != nil || resp.StatusCode >= 400 {
					fmt.Println("failed to download", url)
					continue
				}

				contentType := resp.Header.Get("Content-Type")
				extensions, errr := mime.ExtensionsByType(contentType)
				if errr != nil {
					return errr
				}
				if len(extensions) > 0 {
					fn += extensions[0]
					fp += extensions[0]
				} else {
					fileExtensionRegex := regexp.MustCompile("(\\.[0-9a-z]+)(?:[?#%&]|$)")
					ms := fileExtensionRegex.FindAllStringSubmatch(url, -1)
					if len(ms) > 0 && len(ms[0]) > 1 {
						fn += ms[0][1]
						fp += ms[0][1]
					}
				}

				out, _ := os.Create(fp)
				_, _ = io.Copy(out, resp.Body)
				_ = out.Close()
				_ = resp.Body.Close()

				fmt.Println("downloaded", contentType, url)
				fmt.Println(fp)
			} else {
				return er
			}
		} else {
			fmt.Println("already exists", url)
		}

		str = strings.Replace(str, url, "/"+strings.Split(htmlFn, string(os.PathSeparator))[1]+"/_external/"+fn, 1)
	}

	return os.WriteFile(filepath.Join(htmlDir, htmlFn), []byte(str), os.ModePerm)
}

var blockedStrings = []string{".htm", ".exe", ".pdf", ".zip", ".apk", ".STEP"}

func isBlocked(s string) bool {
	for bs := range blockedStrings {
		if strings.Contains(s, blockedStrings[bs]) {
			return true
		}
	}
	return false
}
