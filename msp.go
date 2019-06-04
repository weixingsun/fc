package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"./serial"
	"log"
	"os"
	"net"
	"bufio"
	"time"
	//"math/rand"
)

const (
	msp_API_VERSION = 1
	msp_FC_VARIANT  = 2
	msp_FC_VERSION  = 3
	msp_BOARD_INFO  = 4
	msp_BUILD_INFO  = 5

	msp_NAME        = 10

	msp_SET_RAW_RC = 200
	msp_RC         = 105

	rx_START = 1400
	rx_RAND  =  200

	state_INIT = iota
	state_M
	state_DIRN
	state_LEN
	state_CMD
	state_DATA
	state_CRC

	LEVEL_AUX = 1500
)

type MSPSerial struct {
	klass int
	p *serial.Port
	conn net.Conn
	reader *bufio.Reader
}

func encode_msp(cmd byte, payload []byte) []byte {
	var paylen byte
	if len(payload) > 0 {
		paylen = byte(len(payload))
	}
	buf := make([]byte, 6+paylen)
	buf[0] = '$'
	buf[1] = 'M'
	buf[2] = '<'
	buf[3] = paylen
	buf[4] = cmd
	if paylen > 0 {
		copy(buf[5:], payload)
	}
	crc := byte(0)
	for _, b := range buf[3:] {
		crc ^= b
	}
	buf[5+paylen] = crc
	return buf
}
func (m *MSPSerial) read(inp []byte) (int, error) {
	if m.klass == DevClass_SERIAL {
		return m.p.Read(inp)
	} else if m.klass == DevClass_TCP {
		return m.conn.Read(inp)
	} else {
		return m.reader.Read(inp)
	}
}

func (m *MSPSerial) write(payload []byte) (int, error) {
	if m.klass == DevClass_SERIAL {
		return m.p.Write(payload)
	} else {
		return m.conn.Write(payload)
	}
}

func (m *MSPSerial) close() {
    m.p.Close()
}
func (m *MSPSerial) Read_msp() (byte, []byte, error) {
	inp := make([]byte, 1)
	var count = byte(0)
	var len = byte(0)
	var crc = byte(0)
	var cmd = byte(0)
	ok := true
	done := false
	var buf []byte

	n := state_INIT

	for !done {
		_, err := m.read(inp)
		if err == nil {
			switch n {
			case state_INIT:
				if inp[0] == '$' {
					n = state_M
				}
			case state_M:
				if inp[0] == 'M' {
					n = state_DIRN
				} else {
					n = state_INIT
				}
			case state_DIRN:
				if inp[0] == '!' {
					n = state_LEN
					ok = false
				} else if inp[0] == '>' {
					n = state_LEN
				} else {
					n = state_INIT
				}
			case state_LEN:
				len = inp[0]
				buf = make([]byte, len)
				crc = len
				n = state_CMD
			case state_CMD:
				cmd = inp[0]
				crc ^= cmd
				if len == 0 {
					n = state_CRC
				} else {
					n = state_DATA
				}
			case state_DATA:
				buf[count] = inp[0]
				crc ^= inp[0]
				count++
				if count == len {
					n = state_CRC
				}
			case state_CRC:
				ccrc := inp[0]
				if crc != ccrc {
					ok = false
				}
				done = true
			}
		}
	}
	if !ok {
		return 0, nil, errors.New("MSP error")
	} else {
		return cmd, buf, nil
	}
}

func (m *MSPSerial) Read_cmd(cmd byte) (byte, []byte, error) {
	var buf []byte
	var err error
	var c = byte(0)

	for ; c != cmd && err ==nil ; {
		c,buf,err = m.Read_msp()
		if c != cmd {
			fmt.Printf("Unsolicted %v (wanted %v)\n", c, cmd)
		}
	}
	return c,buf,err
}

func NewMSPSerial(dd DevDescription) *MSPSerial {
	c := &serial.Config{Name: dd.name, Baud: dd.param}
	p, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	return &MSPSerial{klass: dd.klass, p: p}
}
/////////////////////////////////////////////////////////
func NewMSPTCP(dd DevDescription) *MSPSerial {
	var conn net.Conn
	remote := fmt.Sprintf("%s:%d", dd.name, dd.param)
	addr, err := net.ResolveTCPAddr("tcp", remote)
	if err == nil {
    conn, err = net.DialTCP("tcp", nil, addr)
	}

	if err != nil {
		log.Fatal(err)
	}
	return &MSPSerial{klass: dd.klass, conn: conn}
}

func NewMSPUDP(dd DevDescription) *MSPSerial {
	var laddr, raddr *net.UDPAddr
	var reader  *bufio.Reader
	var conn net.Conn
	var err error

	if dd.param1 != 0 {
		raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dd.name1, dd.param1))
		laddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dd.name, dd.param))
	} else {
		if dd.name == "" {
			laddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dd.name, dd.param))
		} else {
			raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dd.name, dd.param))
		}
	}
	if err == nil {
		conn,err = net.DialUDP("udp", laddr, raddr)
		if err == nil {
		reader = bufio.NewReader(conn)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return &MSPSerial{klass: dd.klass, conn: conn, reader : reader}
}
///////////////////////////////////////////////////////////////////////////
func (m *MSPSerial) Send_msp(cmd byte, payload []byte) {
	buf := encode_msp(cmd, payload)
	m.write(buf)
}

func MSPInit(dd DevDescription) *MSPSerial {
	var fw, api, vers, board, gitrev string
	var m *MSPSerial
	switch dd.klass {
		case DevClass_SERIAL:
		m = NewMSPSerial(dd)
	case DevClass_TCP:
		m = NewMSPTCP(dd)
	case DevClass_UDP:
		m = NewMSPUDP(dd)
	default:
		fmt.Fprintln(os.Stderr, "Unsupported device")
		os.Exit(1)
	}

	m.Send_msp(msp_API_VERSION, nil)
	_, payload, err := m.Read_cmd(msp_API_VERSION)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: ", err)
	} else {
		api = fmt.Sprintf("%d.%d", payload[1], payload[2])
	}

	m.Send_msp(msp_FC_VARIANT, nil)
	_, payload, err = m.Read_cmd(msp_FC_VARIANT)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: ", err)
	} else {
		fw = string(payload[0:4])
	}

	m.Send_msp(msp_FC_VERSION, nil)
	_, payload, err = m.Read_cmd(msp_FC_VERSION)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: ", err)
	} else {
		vers = fmt.Sprintf("%d.%d.%d", payload[0], payload[1], payload[2])
	}

	m.Send_msp(msp_BUILD_INFO, nil)
	_, payload, err = m.Read_cmd(msp_BUILD_INFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: ", err)
	} else {
		gitrev = string(payload[19:])
	}

	m.Send_msp(msp_BOARD_INFO, nil)
	_, payload, err = m.Read_cmd(msp_BOARD_INFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: ", err)
	} else {
		board = string(payload[9:])
	}

	fmt.Fprintf(os.Stderr, "%s v%s %s (%s) API %s", fw, vers, board, gitrev, api)

	m.Send_msp(msp_NAME, nil)
	_, payload, err = m.Read_cmd(msp_NAME)

	if len(payload) > 0 {
		fmt.Fprintf(os.Stderr, " \"%s\"\n", payload)
	} else {
		fmt.Fprintln(os.Stderr, "")
	}
	return m
}
///////////////////////////////////////////////////////////////////////
func (m *MSPSerial) send_cmds(cmds []uint16) {
	//fmt.Println("MSP.send_cmds:",cmds)
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(cmds[0]))
	binary.LittleEndian.PutUint16(buf[2:4], uint16(cmds[1]))
	binary.LittleEndian.PutUint16(buf[4:6], uint16(cmds[2]))
	binary.LittleEndian.PutUint16(buf[6:8], uint16(cmds[3]))
	binary.LittleEndian.PutUint16(buf[8:10],  uint16(cmds[4]))
	binary.LittleEndian.PutUint16(buf[10:12], uint16(cmds[5]))
	binary.LittleEndian.PutUint16(buf[12:14], uint16(cmds[6]))
	binary.LittleEndian.PutUint16(buf[14:16], uint16(cmds[7]))
	m.Send_msp(msp_SET_RAW_RC, buf)
}

func (m *MSPSerial) arm() {
	DEFAULT_CMD := []uint16 {1500, 1500, 1500, 1000, 1000, 1000, 1000, 1000}
	ARM_CMD     := []uint16 {1500, 1500, 2000, 1000, 1000, 1000, 1000, 1000}
	cnt := 0
	for cnt < 50 {
		m.send_cmds(DEFAULT_CMD)
		time.Sleep(100 * time.Millisecond)
		cnt++
	}
	cnt = 0
	for cnt < 11 {
		m.send_cmds(ARM_CMD)
		time.Sleep(100 * time.Millisecond)
		cnt++
	}
}

func (m *MSPSerial) disarm() {
	DISARM_CMD := []uint16 {1500, 1500, 1000, 1000, 1000, 1000, 1000, 1000}
	cnt := 0
	for cnt < 10 {
		m.send_cmds(DISARM_CMD)
		time.Sleep(100 * time.Millisecond)
		cnt++
	}
}
func (m *MSPSerial) warm(level_throttle uint16) {
    var cnt uint16 = 0
    step:= uint16((level_throttle-1200)/50)
    for cnt < 50 {
        throttle  := 1200+cnt*step
        //fmt.Println("throttle=",throttle," cnt=",cnt," step=",step)
        WARM_CMD  := []uint16 {1500, 1500, 1500, throttle, LEVEL_AUX, 1000, 1000, 1000}
        m.send_cmds(WARM_CMD)
        time.Sleep(100 * time.Millisecond)
        cnt++
    }
}
func (m *MSPSerial) hover(throttle uint16, seconds uint16) {
    HOVER_CMD := []uint16 {1500, 1500, 1500, throttle, LEVEL_AUX, 1000, 1000, 1000}
    var cnt uint16 = 0
    for cnt < seconds*10 {
        m.send_cmds(HOVER_CMD)
        time.Sleep(100 * time.Millisecond)
        cnt++
    }
}

func (m *MSPSerial) land(level_throttle uint16) {
    var cnt uint16 = 0
    step:= uint16((level_throttle-1300)/30)
    for cnt < 30 {
        throttle  := level_throttle-cnt*step
        LAND_CMD  := []uint16 {1500, 1500, 1500, throttle, LEVEL_AUX, 1000, 1000, 1000}
        m.send_cmds(LAND_CMD)
        time.Sleep(100 * time.Millisecond)
        cnt++
    }
}
func (m *MSPSerial) takeoff(level_throttle uint16, hover_time uint16) {
    fmt.Println("arm")
    m.arm()
    fmt.Println("warm")
    m.warm(level_throttle)
    fmt.Println("hover")
    m.hover(level_throttle,hover_time)
    fmt.Println("land")
    m.land(level_throttle)
    fmt.Println("disarm")
    m.disarm()
}
/*
func deserialise_rx(b []byte) ([]int16) {
    buf := make([]int16, 8)
    for j:= 0; j < 8; j++ {
        n := j*2;
        buf[j] = int16(binary.LittleEndian.Uint16(b[n:n+2]))
    }
    return buf
}*/
