package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"net/smtp"
	"os"
	"strings"
	"encoding/base64"
	"path"
)

func main() {
	watchRoot()
	select {}
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
		return
	}
	Watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Println(err)
		return
	}
	go listen(Watcher)
	Watcher.Watch(dir)
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
	err := sendToMail("fltcsong@163.com", sub, content, "html")
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

func listen(self *fsnotify.Watcher) {
	for {
		select {
		case w := <-self.Event:
			log.Println(w)
			if w.IsCreate() {
				if strings.HasSuffix(strings.ToLower(w.Name), ".htm") && strings.Contains(strings.ToLower(w.Name), "output") {
					log.Println("化合物生成结束")
					SendEmail(w.Name)
					self.RemoveWatch(path.Dir(w.Name))
					break
				}

				if isDir(w.Name) {
					log.Println("add watch on ", w.Name, " ...")
					self.Watch(w.Name)
				}
			}
		case err := <-self.Error:
			log.Fatalln(err)
		}
	}
}