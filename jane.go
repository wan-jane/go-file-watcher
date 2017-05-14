package main

import (
    "github.com/howeyc/fsnotify"
    "log"
    "os"
    "strings"
    "fmt"
    "net/smtp"
)

type monitor struct {
    watch *fsnotify.Watcher
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("缺少参数，请 这样使用：jane.exe 目录 邮箱")
        fmt.Println("no file/dirtory param to watch")
        os.Exit(1)
    }
    fileName := os.Args[1]
    _,err := os.Stat(fileName)
    if err != nil {
        log.Println("open stream failed")
        os.Exit(1)
    }
    M, err := NewMonitor()
    if err != nil {
        log.Println(err)
        return
    }
    M.Do()
    M.watch.Watch(fileName)
    select {}
}

func NewMonitor() (monitor, error) {
    Mon, err := fsnotify.NewWatcher()
    return monitor{Mon}, err
}

func fileSize(fileName string) int64 {
    fi, err := os.Stat(fileName)

    if err != nil {
        if os.IsNotExist(err) {
            log.Println("file not exist")
        }

        return 0
    }

    return fi.Size()
}

func sendToMail(to, subject, body, mailtype string) error {
      host := "smtp.qq.com:25"
      auth := smtp.PlainAuth("", "1312980813@qq.com", "vrbmhtvhuiavbach", "smtp.qq.com")
      var content_type string
      if mailtype == "html" {
            content_type = "Content-Type:text/" + mailtype + ";charset=UTF-8"
      } else {
            content_type = "Content-Type:text/plain" + ";charset=UTF-8"
      }

      msg := []byte("Subject:"+ subject + "\r\n" + content_type + "\r\n\r\n" + body)
      err := smtp.SendMail(host, auth, "1312980813@qq.com", []string{to}, msg)
      return err
}

func SendEmail() {
      sub := "亲爱的，化合物分析结束了"
      content := os.Args[1] + "目录下的化合物跑完了"
      if os.Args[2] == "no" {
          return
      }
      err := sendToMail(os.Args[2], sub, content, "html")
      if err != nil {
          log.Println("send email failed")
      }
}

func (self monitor) Do() {
    go func() {
        for {
            select {
            case w := <-self.watch.Event:
                if w.IsModify() {
                    log.Println(w.Name, "is changed, filesize is:",fileSize(w.Name))
                    continue
                }
                if w.IsDelete() {
                    log.Println(w.Name, "is deleted")
                    continue
                }
                if w.IsRename() {
                    log.Println(w.Name, "is renamed")
                    self.watch.RemoveWatch(w.Name)
                    continue
                }
                if w.IsCreate() {
                    log.Println(w.Name, "is created filesize is:", fileSize(w.Name))
                    if strings.HasSuffix(w.Name, ".html") {
                       log.Println("化合物生成结束")
                       SendEmail()
                    }
                    continue
                }
            case err := <-self.watch.Error:
                log.Fatalln(err)
            }
        }
    }()
}
