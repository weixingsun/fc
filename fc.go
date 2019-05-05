package main

import (
	//"fmt"
	"os"
	"log"
	"strings"
	"path/filepath"
	//"flag"
	//"regexp"
	//"strconv"
)
type FC struct {
	cfg *Cfg
	lora *Lora
	msp  *MSPSerial
	cmds map[string]f_str
        level_throttle uint16
        hover_time uint16
}

func NewFC() *FC {
    fc := new(FC)
    fc = fc.initMsg()
    fc.msp  = fc.initMSP()
    fc.cfg  = fc.initCFG()
    fc.lora = fc.initRF()
    return fc
}
func (fc *FC) initMSP() *MSPSerial{
    devdesc := DevDescription{klass: DevClass_SERIAL }
    devdesc.name = "/dev/UART_CF"
    devdesc.param = 115200
    //fmt.Println("MSP: ",devdesc)
    return MSPInit(devdesc)
}
func (fc *FC) initRF() *Lora {
    return NewLora(
        "/dev/UART_LORA",
        9600,
        "utf-8",
        "FFFF",
    )
}
func (fc *FC) initCFG() *Cfg {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
       log.Fatal(err)
    }
    p := dir+"/fc.ini"
    c := NewCfg(p)
    fc.level_throttle = c.geti("level")
    fc.hover_time     = c.geti("hover")
    return c
}
func (fc *FC) initMsg() *FC {
    fc.cmds = map[string]f_str{
        "takeoff": fc.takeoff,
        "land":    nil,
        "stop":    nil,
        "ip":      nil,
        "level":   fc.level,
        "hover":   fc.hover,
        "shutdown":nil,
    }
    return fc
}
func (fc *FC) proc_cmd(cmd string) {
    cmd = strings.TrimRight(cmd, "\r\n")
    f := cmd
    p := ""
    if strings.Contains(cmd, " ") {
        arr := strings.Fields(cmd)
        f=arr[0]
        p=arr[1]
    }
    fn,found := fc.cmds[f]
    if found {
        //fmt.Println("fc.proc_cmd found: f="+f+ " p="+ p)
        fn(p)
    } else {
        //fmt.Println("fc.proc_cmd not found: f="+f+ " p="+ p)
        fc.unknown("")
    }
}
func (fc *FC) unknown(s string) {
    fc.lora.send("unknown cmd")
}
func (fc *FC) takeoff(s string) {
    fc.lora.send("takeoff")
    fc.msp.takeoff(fc.level_throttle, fc.hover_time)
}
func (fc *FC) level(l string) {
    fc.cfg.seta("level",l)
    fc.lora.send("level "+l)
}
func (fc *FC) hover(h string) {
    fc.cfg.seta("hover",h)
    fc.lora.send("hover "+h)
}
func main() {
    f := NewFC()
    f.lora.send("ready to flight")
    f.lora.listen(f.proc_cmd);
    defer f.lora.close()
    defer f.msp.close()
    //defer f.cfg.Close()
}
