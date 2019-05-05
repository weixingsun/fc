package main
import (
    "fmt"
    "io"
    "log"
    "time"
    "./serial"
    //"flag"
    //"strings"
)
//var sbuf string
type Lora struct {
    dev, encoder, addr string
    baudrate int
    port     io.ReadWriteCloser
    //sbuf     string
}

func NewLora(dev string,baud int,addr string,encoder string)*Lora{
    lora := new(Lora)
    lora.dev = dev
    lora.baudrate = baud
    lora.addr = addr
    lora.encoder = encoder
    cfg := &serial.Config{
        Name: dev,
        Baud: baud,
        ReadTimeout: time.Second * 5,
    }
    port, err := serial.OpenPort(cfg)
    if err != nil {
        log.Fatalf("Error Serial.Open: %v", err)
    }
    lora.port = port
    return lora
}

func (l *Lora) close() {
    l.port.Close()
}
func (l *Lora) send(msg string) {
    b := []byte(msg+"\n")
    l.port.Write(b)
    //n,err := l.port.Write(b)
    //if err != nil {
    //    fmt.Println("Error sending from serial port: ", err)
    //}
    //fmt.Printf("send %v in %v bytes", msg, n)
}

func (l *Lora) listen(f f_str) {
    sbuf := ""
    for {
        buf := make([]byte, 32)
        n, err := l.port.Read(buf)
        if err != nil {
            if err != io.EOF {
                fmt.Println("Error reading from serial port: ", err, io.EOF)
            }
        } else {
            buf = buf[:n]
            a := string(buf)
	    if n > 1 {
                //fmt.Println("recv(",  n, ") ",a)
		f(a)
		sbuf=""
            }else if a == "\n" {
                //fmt.Println("recv :",  sbuf)
                f(sbuf)
                sbuf=""
            } else {
                //fmt.Println("sbuf=[", sbuf,"] a=", a)
                sbuf += a
            }
        }
    }
}

func (l Lora) sprint(msg string) {
    fmt.Print(msg)
}
/*
func main() {
    dev  := flag.String("d", "/dev/ttyUSB0", "serial device")
    baud := flag.Int("b", 115200, "serial device")
    mode := flag.String("m", "recv", "recv/send")
    msg := flag.String("s", "Hello World", "msg to send")
    flag.Parse()
    l := NewLora(
        *dev,
        *baud,
        "utf-8",
        "FFFF",
    )
    if strings.Compare(*mode, "send") == 0 {
        l.send(*msg)
    }else{
        l.listen(l.sprint)
    }
    //defer lora.port.Close()  //when to close ?
}
*/
