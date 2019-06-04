package main
import (
        "os"
        "log"
        "path/filepath"
)
type f_str func(string)
const (
        DevClass_NONE = iota
        DevClass_SERIAL
        DevClass_TCP
        DevClass_UDP
)

type DevDescription struct {
        klass int
        name string
        param int
        name1 string
        param1 int
}

func current_dir() string {
    //dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    ex, err := os.Executable()
    if err != nil {
       log.Fatal(err)
    }
    return filepath.Dir(ex)
}
