package main

import (
    "github.com/howeyc/fsnotify"
    "log"
    "os"
)

type monitor struct {
    watch *fsnotify.Watcher
}

func main() {
    if len(os.Args) < 2 {
        log.Println("no dir param to watch")
        os.Exit(1)
    }
    dirName := os.Args[1]
    fi,err := os.Stat(dirName)
    if err != nil {
        log.Println("open stram failed")
        os.Exit(1)
    }

    if !fi.IsDir() {
        log.Println("dir:", dirName , "not exist")
        os.Exit(1)
    }
    M, err := NewMonitor()
    if err != nil {
        log.Println(err)
        return
    }
    M.Do()
    M.watch.Watch(dirName)
    select {}
}

func NewMonitor() (monitor, error) {
    Mon, err := fsnotify.NewWatcher()
    return monitor{Mon}, err
}

func fileSize(fileName string) {
    fi, err := os.Stat(fileName)

    if err != nil {

        if os.IsNotExist(err) {
            log.Println("file not exist")
        }
        return
    }

    log.Println(fi.Size)
}

func (self monitor) Do() {
    go func() {
        for {
            select {
            case w := <-self.watch.Event:
                log.Println(w)
                if w.IsModify() {
                    fileSize(w.Name)
                    continue
                }
                if w.IsDelete() {
                    log.Println("文件", w.Name, "被删除.")
                    continue
                }
                if w.IsRename() {
                    w = <-self.watch.Event
                    log.Println(w)
                    self.watch.RemoveWatch(w.Name)
                    log.Println(w.Name, " 被重命名.")
                }
            case err := <-self.watch.Error:
                log.Fatalln(err)
            }
        }
    }()
}
