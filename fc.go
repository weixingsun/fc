package main

import (
	"fmt"
	"strings"
	//"os"
	//"log"
	//"flag"
	//"regexp"
	//"strconv"
)
type FC struct {
	cfg *Cfg
	lora *Lora
	msp  *MSPSerial
	cmds map[string]f_str
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
    fmt.Println("MSP: ",devdesc)
    return MSPInit(devdesc)
}
func (fc *FC) initRF() *Lora {
    return NewLora(
        //"/dev/UART_CH340G",
        "/dev/UART_LORA",
        9600,
        "utf-8",
        "FFFF",
    )
}
func (fc *FC) initCFG() *Cfg {
    p := "fc.ini"
    return NewCfg(p)
}
func (fc *FC) initMsg() *FC {
    fc.cmds = map[string]f_str{
        "takeoff": fc.takeoff,
        "land":    nil,
        "hover":   nil,
        "stop":    nil,
        "ip":      nil,
        "level":   fc.level,
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
        fmt.Println("fc.proc_cmd: f="+f+ " p="+ p)
        fn(p)
    } else {
        fc.unknown("")
    }
}
func (fc *FC) unknown(s string) {
    fc.lora.send("unknown cmd")
}
func (fc *FC) takeoff(s string) {
    fc.lora.send("takeoff")
    fc.msp.takeoff()
}
func (fc *FC) level(l string) {
    fc.cfg.seta("level",l)
    fc.lora.send("level "+l)
}
func main() {
    f := NewFC()
    f.lora.listen(f.proc_cmd);
    defer f.lora.port.Close()
}
