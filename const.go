package main
import (
	//"bytes"
	"fmt"
        "os"
	"os/exec"
        "log"
        "path/filepath"
	"strings"
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
        cmd_cancel  = "cancel"
        cmd_takeoff = "takeoff"
        cmd_land    = "land"
        cmd_stop    = "stop"
        cmd_ip      = "ip"
        cmd_level   = "level"
        cmd_hover   = "hover"
        cmd_shutdown= "shutdown"
)

const (
	msg_ready    = "ready to flight"
	msg_cancel   = "cancelling current job"
	msg_takeoff  = "taking off"
	msg_unknown  = "unknown cmd"
	msg_shutdown = "shutdown"
)

const (
	str_empty = ""
	str_space = " "
	str_liner = "\r\n"
)

type f_str func(string)
//type RF struct {}
type RF interface {
    listen(f f_str)
    send(string)
    close()
}

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

func os_shutdown() {
    args := strings.Fields("shutdown -h now")
    //args := strings.Fields("ls -lrt /home/sun/jbb/mlc/")
    os_cmd(args)
}
func os_cmd(args []string) {
    out, err := exec.Command(args[0], args[1], args[2]).Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("cmd: %s \n",args)
    fmt.Printf("%s \n",out)
}
func current_dir() string {
    ex, err := os.Executable()
    if err != nil {
       log.Fatal(err)
    }
    return filepath.Dir(ex)
}
func file_exists(path string) bool{
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    } else {
        return true
    }
}
/*func main() {
    os_shutdown()
}*/
