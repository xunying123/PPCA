package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func udp(client net.Conn, proxy [16]string, addr string, n int, nn int, Array []byte, count int) {
	if count != -1 {
		var array [512]byte
		clientAddr, Err := net.ResolveUDPAddr("udp", addr)
		if Err != nil {
			return
		}
		fmt.Println("UDP_addr", addr)
		fmt.Println("clientAddr", clientAddr.IP)
		clientUDP, _ := net.ListenUDP("udp", nil)
		/*if err1 != nil {
			_, _ = client.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
			return
		}*/
		remoteUDP, _ := net.ListenUDP("udp", nil)
		/*if err2 != nil {
			_ = clientUDP.Close()
			_, _ = client.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
			return
		}*/
		/*idDes, portDes, _ := net.SplitHostPort(clientUDP.LocalAddr().String())
		ip_ := net.ParseIP(idDes)
		po_, _ := strconv.Atoi(portDes)
		port := uint16(po_)
		ayp := 0
		if len(ip_) == 16 {
			ayp = 0x04
		} else {
			ayp = 0x01
		}
		res := []byte{0x05, 0x00, 0x00, byte(ayp)}
		res = append(res, ip_...)
		_, _ = client.Write(binary.BigEndian.AppendUint16(res, port))*/

		parent := context.Background()
		ctx, cancel := context.WithCancel(parent)

		defer clientUDP.Close()
		defer remoteUDP.Close()
		go ReceiveFromClient(clientUDP, remoteUDP, clientAddr)
		go ReceiveFromRemote(clientUDP, remoteUDP, clientAddr)
		_, _ = remoteUDP.WriteToUDP(Array[n:n+nn], clientAddr)
		go func() {
			for {
				_, ERr := client.Read(array[:])
				if ERr != nil {
					break
				}
			}
			cancel()
		}()
		select {
		case <-ctx.Done():
		}
		return
	} else {
		dest, _ := net.Dial("tcp", proxy[0])
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
			_, _ = dest.Write([]byte{0x05, 0x03, 0x00})
			N, ERr := io.ReadFull(dest, array[:2])
			if ERr != nil {
				fmt.Println("wrong read")
				return
			}
			if N != 2 || array[0] != 0x05 || array[1] != 0x00 {
				return
			}
		}
		nnn, eee := dest.Read(array[:])
		if eee != nil {
			fmt.Println("wrong read")
			return
		}
		if nnn <= 6 || array[1] != 0x00 {
			fmt.Println("wrong read")
			return
		}
		_, _ = client.Write(array[:nnn])
		_, _ = dest.Write(Array[n : n+nn])
		transfer(dest, client)
	}
}

func ReceiveFromClient(client *net.UDPConn, server *net.UDPConn, add *net.UDPAddr) {
	var Array [512]byte
	for {
		n, Addr, err := client.ReadFromUDP(Array[:])
		if err != nil {
			break
		}
		if (Addr.IP.To16().String() != add.IP.To16().String()) || (add.Port != Addr.Port) {
			continue
		}
		a, b, c, d := Array[0], Array[1], Array[2], Array[3]
		if (a != 0x00) || (b != 0x00) {
			fmt.Println("no rsv")
			continue
		}
		if c != 0x00 {
			fmt.Println("no frag")
			continue
		}
		index := 0
		addr := ""
		switch d {
		case 0x01:
			{
				port := binary.BigEndian.Uint16(Array[8:10])
				addr = fmt.Sprintf("%d.%d.%d.%d:%d", Array[4], Array[5], Array[6], Array[7], port)
				index = 10
				break
			}
		case 0x03:
			{
				a := Array[4]
				port := binary.BigEndian.Uint16(Array[a+6 : a+8])
				addr = string(Array[5:a+6]) + fmt.Sprintf(":%v", port)
				index = int(Array[4]) + 6
				break
			}
		case 0x04:
			{
				port := binary.BigEndian.Uint16(Array[20:22])
				addr = "["
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[4:6]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[6:8]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[8:10]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[10:12]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[12:14]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[14:16]))
				addr += fmt.Sprintf("%x:", binary.BigEndian.Uint16(Array[16:18]))
				addr += fmt.Sprintf("%x", binary.BigEndian.Uint16(Array[18:20]))
				addr += fmt.Sprintf("]:%d", port)
				index = 22
				break
			}
		}
		remote, ERR := net.ResolveUDPAddr("udp", addr)
		if ERR != nil {
			fmt.Println("resolve wrong")
			continue
		}
		_, _ = server.WriteToUDP(Array[index:n], remote)
	}
}

func ReceiveFromRemote(client *net.UDPConn, server *net.UDPConn, addr *net.UDPAddr) {
	var array [512]byte
	for {
		n, _, err := server.ReadFromUDP(array[:])
		if err != nil {
			fmt.Println("wrong receive")
			continue
		}
		ss := []byte{0x00, 0x00, 0x00}
		ipp := addr.IP
		Port := addr.Port
		if ipp.To16() != nil {
			ss = append(ss, 0x01)
			ss = append(ss, binary.BigEndian.AppendUint16(ipp.To4(), uint16(Port))...)
		} else {
			ss = append(ss, 0x04)
			ss = append(ss, binary.BigEndian.AppendUint16(ipp.To4(), uint16(Port))...)
		}
		_, _ = client.WriteToUDP(append(ss[:], array[0:n]...), addr)
	}
}
