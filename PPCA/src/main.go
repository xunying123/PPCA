package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {

	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen failed", err)
	}
	defer server.Close()
	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println("listen failed")
			continue
		}
		go process(client)
	}
}

func process(client net.Conn) {
	fmt.Println("processing")
	array := make([]byte, 32*1024)
	n, Err := io.ReadFull(client, array[:2])
	if Err != nil {
		fmt.Println("wrong read")
		return
	}
	methodNum := int(array[1])
	n, _ = io.ReadFull(client, array[:methodNum])
	if n != methodNum {
		fmt.Println("wrong method num")
		return
	}
	auth := false
	for i := 0; i < methodNum; i++ {
		if array[i] == 0x00 {
			auth = true
			break
		}
	}
	if !auth {
		_, _ = client.Write([]byte{0x05, 0xff})
		fmt.Println("wrong auth")
		return
	}
	_, _ = client.Write([]byte{0x05, 0x00})
	fmt.Println(2)
	_, eee := io.ReadFull(client, array[:4])
	if eee != nil {
		fmt.Println("read error")
		return
	}
	if (array[0] != 0x05) || (array[1] != 0x01 && array[1] != 0x04) || (array[2] != 0x00) {
		_, _ = client.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		fmt.Println("wrong cmd")
		return
	}
	ayp := array[3]
	ype := array[1]
	addr := ""
	support := true
	switch ayp {
	case 0x01:
		{
			_, _ = io.ReadFull(client, array[4:10])
			n = 10
			port := binary.BigEndian.Uint16(array[8:10])
			addr = fmt.Sprintf("%d.%d.%d.%d:%d", array[4], array[5], array[6], array[7], port)
			break
		}
	case 0x03:
		{
			_, _ = io.ReadFull(client, array[4:5])
			a := array[4]
			_, _ = io.ReadFull(client, array[5:7+a])
			n = int(7 + a)
			port := binary.BigEndian.Uint16(array[a+5 : a+7])
			addr = string(array[5:a+5]) + fmt.Sprintf(":%v", port)
			break
		}
	case 0x04:
		{
			_, _ = io.ReadFull(client, array[4:22])
			n = 22
			port := binary.BigEndian.Uint16(array[20:22])
			addr = "["
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[4:6]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[6:8]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[8:10]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[10:12]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[12:14]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[14:16]))
			addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(array[16:18]))
			addr += fmt.Sprintf("%x", binary.BigEndian.Uint16(array[18:20]))
			addr += fmt.Sprintf("]:%d", port)
			break
		}
	default:
		{
			support = false
		}
	}
	if support == false {
		_, _ = client.Write([]byte{0x05, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		fmt.Println("unsupported address")
		return
	}
	fmt.Println(addr)
	nn, _ := client.Read(array[n : n+10240])
	var proxy [16]string
	var count = 0
	proxy, eee, count = http(array)
	if eee == nil {
		if ype == 0x01 {
			tcp(client, proxy, addr, n, nn, array, count)
		} else {
			udp(client, proxy, addr, n, nn, array, count)
		}
	} else {
		proxy, eee, count = tls(array, n)
		if eee == nil {
			if ype == 0x01 {
				tcp(client, proxy, addr, n, nn, array, count)
			} else {
				udp(client, proxy, addr, n, nn, array, count)
			}
		} else {
			proxy, eee, count = pid()
			if eee == nil {
				if ype == 0x01 {
					tcp(client, proxy, addr, n, nn, array, count)
				} else {
					udp(client, proxy, addr, n, nn, array, count)
				}
			} else {
				proxy, count = divide(addr, ayp)
				if ype == 0x01 {
					tcp(client, proxy, addr, n, nn, array, count)
				} else {
					udp(client, proxy, addr, n, nn, array, count)
				}
			}
		}
	}
}
