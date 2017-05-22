package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"net/smtp"
	"os"
	"strings"
	"encoding/base64"
	"errors"
)

var watchPaths []string

type JaneWatcher struct {
	watch *fsnotify.Watcher
}

func main() {
	watchRoot()
	select {}
}

func newJaneWatcher() (*JaneWatcher, error) {
	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return nil, errors.New("can not init !")
	}

	return &JaneWatcher{Watcher}, nil
}

func watchRoot() {
	rootDir := "/Users/xnw/gopath/src/go-file-watcher/test"
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	startWatch(rootDir)
}

func startWatch(dir string)  {
	_, err := os.Stat(dir)
	if err != nil {
		log.Println("open stream " + dir + " failed")
		os.Exit(1)
	}
	janeW, err:= newJaneWatcher()
	if err != nil {
		log.Println(err)
		return
	}
	go janeW.listen()
	janeW.watch.WatchFlags(dir, fsnotify.FSN_CREATE)
}

func sendToMail(to, subject, body, mailtype string) error {
	host := "smtp.qq.com:25"
	auth := smtp.PlainAuth("", "1312980813@qq.com", "vrbmhtvhuiavbach", "smtp.qq.com")
	var content_type string = "MIME-Version: 1.0\r\n"
	subject = "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + ";charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + ";charset=UTF-8"
	}
    content_type = content_type + "Content-Transfer-Encoding: 8bit\r\n"
	msg := []byte("To:" + to + "\r\n"+"Subject:" + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(host, auth, "1312980813@qq.com", []string{to}, msg)
	return err
}

func SendEmail(fileName string) {
	sub := "亲爱的，化合物分析结束了"
	content := "已经生成了" + fileName + "文件"
	email := "fltcsong@163.com"
	if len(os.Args) > 2 {
		email = os.Args[2]
	}
	err := sendToMail(email, sub, content, "html")
	if err != nil {
		log.Println("send email failed")
        	return
	}

    log.Println("邮件已发送")
}

func isDir(pathname string) bool {
	fi, err := os.Stat(pathname)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func (jw *JaneWatcher) listen() {
	for {
		select {
		case w := <-jw.watch.Event:
			log.Println(w)
			if w.IsCreate() {
				if strings.Contains(strings.ToLower(w.Name), "output") && strings.HasSuffix(strings.ToLower(w.Name), "report.htm") {
					log.Println("化合物生成结束")
					SendEmail(w.Name)
					jw.removeAllChildWatch()
					continue
				}

				if isDir(w.Name) {
					log.Println("add watch on ", w.Name, " ...")
					jw.watch.WatchFlags(w.Name, fsnotify.FSN_CREATE)
					watchPaths = append(watchPaths, w.Name)
					printWatchers()
				}
			}
		case err := <-jw.watch.Error:
			log.Fatalln(err)
		}
	}
}

func (jw *JaneWatcher) removeAllChildWatch() {
	for _, path := range watchPaths {
		jw.watch.RemoveWatch(path)
	}

	log.Println("remove all child watch ...")
	watchPaths = []string{}
	printWatchers()
}

func printWatchers() {
	log.Println("--------已监听子目录---------")
	for _, path := range watchPaths {
		log.Println(path)
	}
	log.Println("---------------------------")
}