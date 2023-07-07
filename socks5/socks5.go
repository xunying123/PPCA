package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func main() {
	sever, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen failed", err)
	}
	defer sever.Close()
	for {
		client, err := sever.Accept()
		if err != nil {
			fmt.Println("listen failed")
			continue
		}
		go process(client)
	}
}

func process(client net.Conn) {
	err := socksReceive1(client)
	if err != nil {
		fmt.Println("Receive Wrong", err)
		return
	}
	target, errr := socksReceive2(client)
	if errr != nil {
		fmt.Println("connect failed", err)
		return
	}
	transfer(client, target)
}

func socksReceive1(client net.Conn) (err error) {
	array := make([]byte, 512)
	n, Err := io.ReadFull(client, array[:2])
	if Err != nil {
		return errors.New("reading head Failed")
	}
	methodNum := int(array[1])
	n, _ = io.ReadFull(client, array[:methodNum])
	if n != methodNum {
		return errors.New("wrong methodnum")
	}
	auth := false
	for i := 0; i < methodNum; i++ {
		if array[i] == 0x00 {
			auth = true
			break
		}
	}
	if !auth {
		client.Write([]byte{0x05, 0xff})
		return errors.New("wrong auth")
	}
	client.Write([]byte{0x05, 0x00})
	return nil
}

func socksReceive2(client net.Conn) (server net.Conn, err error) {
	var array [512]byte
	_, eee := io.ReadFull(client, array[:4])
	if eee != nil {
		return nil, errors.New("read error")
	}
	if (array[0] != 0x05) || (array[1] != 0x01) || (array[2] != 0x00) {
		client.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		return nil, errors.New("wrong cmd")
	}
	atyp := array[3]
	addr := ""

	support := true
	switch atyp {
	case 0x01:
		{
			_, er := io.ReadFull(client, array[:6])
			if er != nil {
				return nil, errors.New("wrong ip")
			}
			port := binary.BigEndian.Uint16(array[4:6])
			addr = fmt.Sprintf("%d.%d.%d.%d:%d", array[0], array[1], array[2], array[3], port)
			break
		}
	case 0x03:
		{
			_, er := io.ReadFull(client, array[:1])
			if er != nil {
				return nil, errors.New("domain")
			}
			a := array[0]
			_, er = io.ReadFull(client, array[:a+2])
			if er != nil {
				return nil, errors.New("domain")
			}
			port := binary.BigEndian.Uint16(array[a : a+2])
			addr = string(array[0:a]) + fmt.Sprintf(":%v", port)
			break
		}
	case 0x04:
		{
			_, er := io.ReadFull(client, array[:18])
			if er != nil {
				return nil, errors.New("ipv6")
			}
			port := binary.BigEndian.Uint16(array[16:18])
			addr = "["
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[0:2]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[2:4]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[4:6]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[6:8]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[8:10]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[10:12]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[12:14]))
			addr += fmt.Sprintf("%x", binary.BigEndian.Uint16(array[14:16]))
			addr += fmt.Sprintf("]:%d", port)
			break
		}
	default:
		{
			support = false
		}
	}
	fmt.Println(1)
	if support == false {
		client.Write([]byte{0x05, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		return nil, errors.New("unsupported address")
	}
	fmt.Println(2)
	dest, Err := net.Dial("tcp", addr)
	fmt.Println(Err)
	if Err != nil {
		var failed byte = 0x00
		if strings.Contains(Err.Error(), "invalid.invalid") {
			failed = 0x04
		} else if strings.Contains(Err.Error(), "connection refused") {
			failed = 0x05
		} else if strings.Contains(Err.Error(), "no route") {
			failed = 0x03
		} else if strings.Contains(Err.Error(), "i/o timeout") {
			failed = 0x04
		} else if strings.Contains(Err.Error(), "network is unreachable") {
			failed = 0x03
		} else if strings.Contains(Err.Error(), "failure in name resolution") {
			failed = 0x04
		}
		client.Write([]byte{0x05, failed, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		return nil, errors.New("no connection")

	}
	ip_des, port_des, _ := net.SplitHostPort(dest.LocalAddr().String())
	ip_ := net.ParseIP(ip_des)
	po_, _ := strconv.Atoi(port_des)
	port := uint16(po_)

	if len(ip_) == 16 {
		atyp = 0x04
	} else {
		atyp = 0x01
	}

	res := []byte{0x05, 0x00, 0x00, byte(atyp)}
	res = append(res, ip_...)
	client.Write(binary.BigEndian.AppendUint16([]byte(res), port))
	return dest, nil
}

func transfer(client, target net.Conn) {
	go copying(client, target)
	go copying(target, client)
}

func copying(client, target net.Conn) {
	defer client.Close()
	defer target.Close()
	io.Copy(client, target)
}
