package main
import (
    "fmt"
    "io"
    "log"
    "./serial"
    //"flag"
    //"os"
    //"encoding/hex"
)
type fn func(string)
//type RF interface {
//    send(string)
//    listen(fn)
//}
type Lora struct {
    //send() void
    //listen() void
    dev, encoder, addr string
    baudrate uint
    //oo       serial.OpenOptions
    port     io.ReadWriteCloser
}

func NewLora(dev string,baudrate uint,addr string,encoder string)*Lora{
    lora := new(Lora)
    lora.dev = dev
    lora.baudrate = baudrate
    lora.addr = addr
    lora.encoder = encoder
    oo := serial.OpenOptions{
        PortName: dev,
        BaudRate: baudrate,
        DataBits: 8,
        StopBits: 1,
        MinimumReadSize: 4,
    }
    port, err := serial.Open(oo)
    if err != nil {
        log.Fatalf("Error Serial.Open: %v", err)
    }
    lora.port = port
    return lora
}

func (l *Lora) send(msg string) {
    //txData_, err := hex.DecodeString(msg)
    b := []byte(msg)
    l.port.Write(b)
}

func (l *Lora) listen(f fn) {
    for {
        buf := make([]byte, 32)
        n, err := l.port.Read(buf)
        if err != nil {
            if err != io.EOF {
                fmt.Println("Error reading from serial port: ", err)
            }
        } else {
            buf = buf[:n]
            f(string(buf))
        }
    }
}

func (l Lora) sprint(msg string) {
    fmt.Println(msg)
}
func main() {
    l := NewLora(
        "/dev/UART_CP2102",
        9600,
        "utf-8",
        "FFFF",
    )
    //defer lora.port.Close()  //when to close ?
    l.listen(l.sprint)
    //l.send("haha")
}
