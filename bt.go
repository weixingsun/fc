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
    nfd    int
    buf    []byte
}

func (bt *BT) send(msg string) {
    b := []byte(msg+"\n")
    unix.Write(bt.nfd, b)
}
func (bt *BT) listen(f f_str) {
    for {
        //fmt.Println("reading bt")
        _, err:= unix.Read(bt.nfd, bt.buf)
        if err != nil {
            //log.Println(err)
            fmt.Println("error when reading bt")
            bt.close()
            time.Sleep(2*time.Second)
            bt.fd,bt.nfd = bt.init()
        }else{
            //fmt.Println("Read successfully")
            n0  := bytes.Index(bt.buf, []byte{0})
            n13 := bytes.Index(bt.buf, []byte{13})  //13=\r
            n10 := bytes.Index(bt.buf, []byte{10})  //10=\n
            n := min3(n0,n13,n10)
            s0 := string(bt.buf[:n])
            s := strings.TrimSpace(s0)
            //printAscii(s)
            //fmt.Printf("Received: %v chars: %s\n", len(s), s )
            f(s)
        }
    }
}

func NewBT()*BT{
    bt := new(BT)
    bt.buf = make([]byte, 30)
    bt.fd,bt.nfd  = bt.init()
    return bt
}

func (bt *BT) init() (int,int){
    fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
    if err != nil {
        fmt.Println("Init BT socket err")
        log.Println(err)
    }
    //fmt.Println("Init BT socket sucessfully")
    addr := &unix.SockaddrRFCOMM{
        Channel: 1,
        Addr:    [6]uint8{0,0,0,0,0,0},
    }
    _ = unix.Bind(fd, addr)
    //fmt.Println("Bind BT socket sucessfully")
    _ = unix.Listen(fd,1)
    //fmt.Println("Listen BT socket sucessfully")
    nfd, sa, _ := unix.Accept(fd)
    //print client mac addr
    fmt.Printf("Conn addr=%v fd=%d\n", sa.(*unix.SockaddrRFCOMM).Addr, nfd)
    return fd,nfd
}
func (bt *BT) close() {
    unix.Close(bt.nfd)
    unix.Close(bt.fd)
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
