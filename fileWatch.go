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
        log.Println("no file param to watch")
        os.Exit(1)
    }
    fileName := os.Args[1]
    _,err := os.Stat(fileName)
    if err != nil {
        log.Println("open stram failed")
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
                    continue
                }
            case err := <-self.watch.Error:
                log.Fatalln(err)
            }
        }
    }()
}
