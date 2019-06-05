package main

import (
	//"fmt"
	//"os"
	//"log"
	"strings"
	//"path/filepath"
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
    devdesc.name = DEV_NAME_MSP
    devdesc.param = DEV_BAUDRATE_MSP
    return MSPInit(devdesc)
}
func (fc *FC) initRF() *Lora {
    return NewLora(
        DEV_NAME_LORA,
        DEV_BAUDRATE_LORA,
        DEV_LORA_ENCODER,
        DEV_LORA_ADDR,
    )
}
func (fc *FC) initCFG() *Cfg {
    dir := current_dir()
    p := dir+FILE_CFG
    c := NewCfg(p)
    fc.level_throttle = c.geti(tag_cfg_level_throttle)
    fc.hover_time     = c.geti(tag_cfg_hover_time)
    return c
}
func (fc *FC) initMsg() *FC {
    fc.cmds = map[string]f_str{
        cmd_takeoff: fc.takeoff,
        cmd_land:    nil,
        cmd_stop:    nil,
        cmd_ip:      nil,
        cmd_level:   fc.level,
        cmd_hover:   fc.hover,
        cmd_shutdown:nil,
    }
    return fc
}
func (fc *FC) proc_cmd(cmd string) {
    cmd = strings.TrimRight(cmd, str_liner)
    f := cmd
    p := str_empty
    if strings.Contains(cmd, str_space) {
        arr := strings.Fields(cmd)
        f=arr[0]
        p=arr[1]
    }
    fn,found := fc.cmds[f]
    if found {
        fn(p)
    } else {
        fc.unknown(str_empty)
    }
}
func (fc *FC) unknown(s string) {
    fc.lora.send(msg_unknown)
}
func (fc *FC) takeoff(s string) {
    fc.lora.send(msg_takeoff)
    fc.msp.takeoff(fc.level_throttle, fc.hover_time)
}
func (fc *FC) level(l string) {
    fc.cfg.seta(tag_cfg_level_throttle,l)
    fc.lora.send(tag_cfg_level_throttle+str_space+l)
}
func (fc *FC) hover(h string) {
    fc.cfg.seta(tag_cfg_hover_time,h)
    fc.lora.send(tag_cfg_hover_time+str_space+h)
}
func main() {
    f := NewFC()
    f.lora.send(msg_ready)
    f.lora.listen(f.proc_cmd);
    defer f.lora.close()
    defer f.msp.close()
}
