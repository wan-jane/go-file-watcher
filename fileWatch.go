package main

import (
    "github.com/howeyc/fsnotify"
    "log"
)

type monitor struct {
    watch *fsnotify.Watcher
}

func main() {
    M, err := NewMonitor()
    if err != nil {
        log.Println(err)
        return
    }
    M.Do()
    M.watch.Watch("/ssd/gengkang/.bashrc")
    select {}
}

func NewMonitor() (monitor, error) {
    Mon, err := fsnotify.NewWatcher()
    return monitor{Mon}, err
}

func (self monitor) Do() {
    go func() {
        for {
            select {
            case w := <-self.watch.Event:
                log.Println(w)
                if w.IsModify() {
                    continue
                }
                if w.IsDelete() {
                    log.Println("文件被删除.")
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
