package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func tcp(client net.Conn, proxy [16]string, addr string, n int, nn int, Array []byte, count int) {
	if count == 0 {
		dest, _ := net.Dial("tcp", addr)
		/*fmt.Println(Err)
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
			_, _ = client.Write([]byte{0x05, failed, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
			return
		}
		ipDes, portDes, _ := net.SplitHostPort(dest.LocalAddr().String())
		ip_ := net.ParseIP(ipDes)
		po_, _ := strconv.Atoi(portDes)
		port := uint16(po_)
		ayp := 0x00
		if len(ip_) == 16 {
			ayp = 0x04
		} else {
			ayp = 0x01
		}

		res := []byte{0x05, 0x00, 0x00, byte(ayp)}
		res = append(res, ip_...)
		_, _ = client.Write(binary.BigEndian.AppendUint16(res, port))*/
		_, _ = dest.Write(Array[n : n+nn])
		transfer(client, dest)
	} else {
		dest, _ := net.Dial("tcp", proxy[0])
		fmt.Println("tcp start")
		/*if Err != nil {
			var failed byte = 0x00
			if strings.Contains(Err.Error(), "proxy invalid.invalid") {
				failed = 0x04
			} else if strings.Contains(Err.Error(), "proxy connection refused") {
				failed = 0x05
			} else if strings.Contains(Err.Error(), "proxy no route") {
				failed = 0x03
			} else if strings.Contains(Err.Error(), "proxy i/o timeout") {
				failed = 0x04
			} else if strings.Contains(Err.Error(), "proxy network is unreachable") {
				failed = 0x03
			} else if strings.Contains(Err.Error(), "proxy failure in name resolution") {
				failed = 0x04
			}
			_, _ = client.Write([]byte{0x05, failed, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
			return
		}*/
		_, _ = dest.Write([]byte{0x05, 0x01, 0x00})
		array := make([]byte, 32*1024)
		N, ERr := io.ReadFull(dest, array[:2])
		if ERr != nil {
			fmt.Println("wrong read")
			return
		}
		if N != 2 || array[0] != 0x05 || array[1] != 0x00 {
			return
		}
		for i := 1; i < count; i++ {
			fmt.Println("for")
			ayp := 0
			ipDes, portDes, _ := net.SplitHostPort(proxy[i])
			ip_ := net.ParseIP(ipDes)
			po_, _ := strconv.Atoi(portDes)
			port := uint16(po_)
			if len(ip_) == 16 {
				ayp = 0x04
			} else {
				ayp = 0x01
			}
			res := []byte{0x05, 0x01, 0x00, byte(ayp)}
			res = append(res, ip_...)
			_, _ = dest.Write(binary.BigEndian.AppendUint16(res, port))
			nnn, eee := dest.Read(array[:])
			if eee != nil {
				fmt.Println("wrong read")
				return
			}
			if nnn <= 6 || array[1] != 0x00 {
				fmt.Println("wrong read")
				return
			}
			_, _ = dest.Write([]byte{0x05, 0x01, 0x00})
			N, ERr := io.ReadFull(dest, array[:2])
			if ERr != nil {
				fmt.Println("wrong read")
				return
			}
			if N != 2 || array[0] != 0x05 || array[1] != 0x00 {
				return
			}
		}
		fmt.Println("write")
		_, _ = dest.Write(Array[:n])
		nnn, eee := dest.Read(array[:])
		if eee != nil {
			fmt.Println("wrong read")
			return
		}
		if nnn <= 6 || array[1] != 0x00 {
			fmt.Println("wrong read")
			return
		}
		_, _ = dest.Write(Array[n : n+nn])
		transfer(client, dest)
	}
}

func transfer(client, target net.Conn) {
	go copying(client, target)
	go copying(target, client)
}

func copying(client, target net.Conn) {
	defer client.Close()
	defer target.Close()
	_, _ = io.Copy(client, target)
}
