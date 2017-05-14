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

func SendEmail() {
      sub := "亲爱的，化合物跑完了"
      content := os.Args[1] + "目录下的化合物跑完了"
      mailList := []string{os.Args[2]}
      auth := smtp.PlainAuth(
              "",
              "1312980813@qq.com",
              "21131P13njcit",
              "smtp.qq.com",
              //"smtp.gmail.com",
        )
        // Connect to the server, authenticate, set the sender and recipient,
        // and send the email all in one step.
        err := smtp.SendMail(
                "smtp.qq.com:25",
                auth,
                "1312980813@qq.com",
                mailList,
                []byte(sub+content),
        )
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
                    if strings.HasSuffix(w.Name, ".html") {
                       log.Println("化合物生成结束")
                       SendEmail()
                    }
                    log.Println(w.Name, "is created filesize is:", fileSize(w.Name))
                    continue
                }
            case err := <-self.watch.Error:
                log.Fatalln(err)
            }
        }
    }()
}
