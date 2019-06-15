package main
import (
	//"flag"
	"fmt"
	"log"
	"bytes"
	"strings"
	"syscall"
	"time"
	"golang.org/x/sys/unix"
)

type BT struct {
    RF
    fd     int
    buf    []byte
}
//type f_str func(string)
func min2(x, y int) int {
    if x < y {
        return x
    }
    return y
}
func min3(x, y, z int) int {
    a := min2(x,y)
    if a < z {
        return a
    }
    return z
}

func printAscii(s string) {
    for i,c := range s {
        fmt.Printf("[%d]%c(%d)\n",i,c,int(c))
    }
    fmt.Println()
}

func (bt *BT) send(msg string) {
    b := []byte(msg+"\n")
    unix.Write(bt.fd, b)
}
func (bt *BT) listen(f f_str) {
    for {
        _, err:= unix.Read(bt.fd, bt.buf)
        if err != nil {
            //log.Println(err)
            fmt.Println("error when reading bt")
            time.Sleep(2*time.Second)
            bt.fd = bt.init()
        }else{
            fmt.Println("Read successfully")
            n0  := bytes.Index(bt.buf, []byte{0})
            n13 := bytes.Index(bt.buf, []byte{13})  //13=\r
            n10 := bytes.Index(bt.buf, []byte{10})  //10=\n
            n := min3(n0,n13,n10)
            s0 := string(bt.buf[:n])
            s := strings.TrimSpace(s0)
            //printAscii(s)
            fmt.Printf("Received: %v chars: %s\n", len(s), s )
            f(s)
        }
    }
}

func NewBT()*BT{
    bt := new(BT)
    bt.buf = make([]byte, 30)
    bt.fd  = bt.init()
    return bt
}

func (bt *BT) init() int{
    fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
    if err != nil {
        fmt.Println("Init BT socket err")
        log.Println(err)
    }
    fmt.Println("Init BT socket sucessfully")
    addr := &unix.SockaddrRFCOMM{
        Channel: 1,
        Addr:    [6]uint8{0,0,0,0,0,0},
    }
    _ = unix.Bind(fd, addr)
    fmt.Println("Bind BT socket sucessfully")
    _ = unix.Listen(fd,1)
    fmt.Println("Listen BT socket sucessfully")
    nfd, sa, err := unix.Accept(fd)
    if err != nil {
        fmt.Printf("Waiting... fd=%v  sa=%v", nfd, sa)
        log.Println(err)
        time.Sleep(5*time.Second)
        return 0
    }else{
        //print client mac addr
        //fmt.Printf("Conn addr=%v fd=%d\n", sa.(*unix.SockaddrRFCOMM).Addr, nfd)
        return nfd
    }
}
func (bt *BT) close() {
    //l.port.Close()
}

func (bt *BT) sprint(msg string) {
    fmt.Print(msg)
}
/*
func main() {
    mode := flag.String("m", "recv", "recv/send")
    msg := flag.String("s", "takeoff", "cmd to send")
    flag.Parse()
    bt := NewBT()
    if strings.Compare(*mode, "send") == 0 {
        bt.send(*msg)
    }else{
        bt.listen(bt.send)
    }
}*/
