package main
import (
    //"fmt"
    "log"
    "os"
    "strconv"
    "./ini"
)
type Cfg struct {
    path string
    c    *ini.File
}

func NewCfg(path string) *Cfg{
    cfg := new(Cfg)
    cfg.path = path
    c, err := ini.Load(path)
    if err != nil {
        log.Fatalf("Error Ini.Open: %v", err)
	os.Exit(1)
    }
    cfg.c = c
    return cfg
}

func (c *Cfg) geta(k string) string {
    v := c.c.Section("").Key(k).String()
    //fmt.Println("ini.get( %v ) = %v ", k, v)
    return v
}

func (c *Cfg) geti(k string) uint16 {
    vs := c.geta(k)
    vi,err := strconv.Atoi(vs)
    if err != nil {
        log.Fatalf("Error Ini.Geti: %v", err)
	os.Exit(1)
    }
    return uint16(vi)
}

func (c *Cfg) seta(k string,v string) {
    c.c.Section("").Key(k).SetValue(v)
    c.c.SaveTo(c.path)
}

func (c *Cfg) seti(k string,v int) {
    s:=strconv.Itoa(v)
    c.c.Section("").Key(k).SetValue(s)
    c.c.SaveTo(c.path)
}
/*func main() {
    c := NewCfg(
        "cfg.ini",
    )
    c.seti("level",1400)
    fmt.Printf("ini.set( level ) = 1400 " )
    fmt.Printf("ini.level=",c.get("level"))
}*/
