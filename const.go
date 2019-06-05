package main
import (
        "os"
        "log"
        "path/filepath"
)
const (
	DEV_NAME_MSP  = "/dev/UART_CF"
	DEV_NAME_LORA = "/dev/UART_LORA"
	DEV_BAUDRATE_MSP  = 115200
	DEV_BAUDRATE_LORA = 9600
	DEV_LORA_ENCODER = "utf-8"
	DEV_LORA_ADDR = "FFFF"
)
const (
	FILE_CFG               ="/fc.ini"
	tag_cfg_level_throttle = "level"
	tag_cfg_hover_time     = "hover"
)

const (
	cmd_takeoff = "takeoff"
        cmd_land    = "land"
        cmd_stop    = "stop"
        cmd_ip      = "ip"
        cmd_level   = "level"
        cmd_hover   = "hover"
        cmd_shutdown= "shutdown"
)

const (
	msg_ready = "ready to flight"
	msg_takeoff = "taking off"
	msg_unknown = "unknown cmd"
)

const (
	str_empty = ""
	str_space = " "
	str_liner = "\r\n"
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
