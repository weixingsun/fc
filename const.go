package main
/*import (
	"flag"
)*/
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

/*var (
        baud   = flag.Int("b", 115200, "Baud rate")
        device = flag.String("d", "", "Serial Device")
        arm = flag.Bool("a", false, "Arm (take care now)")
)*/
