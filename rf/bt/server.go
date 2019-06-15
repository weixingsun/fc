package main

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"syscall"
	"time"
	"golang.org/x/sys/unix"
)
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

func init_socket() int {
    fd, _ := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
    addr := &unix.SockaddrRFCOMM{
        Channel: 1,
        Addr:    [6]uint8{0,0,0,0,0,0},
    }
    _ = unix.Bind(fd, addr)
    _ = unix.Listen(fd,1)
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
func main() {
    buf := make([]byte, 30)
bt_init:
    nfd := init_socket()
    if nfd == 0 {
        goto bt_init
    }
    for {
        _, err:= unix.Read(nfd, buf)
        if err != nil {
            //log.Println(err)
            fmt.Println("error when reading bt")
            time.Sleep(2*time.Second)
            goto bt_init
        }
        fmt.Println("Connected")
        n0 := bytes.Index(buf, []byte{0})
        n13 := bytes.Index(buf, []byte{13})  //13=\r
        n10 := bytes.Index(buf, []byte{10})  //10=\n
        n := min3(n0,n13,n10)
        s0 := string(buf[:n])
        s := strings.TrimSpace(s0)
        //printAscii(s)
        fmt.Printf("Received: %v chars: %s\n", len(s), s )
    }
}
