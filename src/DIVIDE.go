package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func divide(addr string, a byte) ([16]string, int) {
	var proxy [16]string
	var count = -1
	switch a {
	case 0x01:
		{
			file, err := os.Open("ipv4.txt")
			if err != nil {
				fmt.Println("无法打开文件:", err)
				return proxy, count
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line[0] == '0' {
					n := len(line)
					var cidr = line[2:n]
					ip, ipNet, _ := net.ParseCIDR(cidr)
					mask := ipNet.Mask.String()
					var ll = ip.String()
					lastIP := make(net.IP, len(ip))
					for i := range ip {
						lastIP[i] = ip[i] | ^mask[i]
					}
					if count > 0 {
						return proxy, count
					}
					var rr = lastIP.String()

					l := uint32(ll[0])<<24 | uint32(ll[1])<<16 | uint32(ll[2])<<8 | uint32(ll[3])
					r := uint32(rr[0])<<24 | uint32(rr[1])<<16 | uint32(rr[2])<<8 | uint32(rr[3])
					a := uint32(addr[0])<<24 | uint32(addr[1])<<16 | uint32(addr[2])<<8 | uint32(addr[3])

					if (a > l) && (a < r) {
						count = 0
					} else {
						count = -1
					}
				} else {
					if count >= 0 {
						n := len(line)
						proxy[count] = line[2:n]
						count++
					}
				}

			}
			return proxy, count
		}
	case 0x03:
		{
			/*pacURL := "http://example.com/proxy.pac"
			u, err := url.Parse(pacURL)
			if err != nil {
				fmt.Println("无法获取 PAC:", err)
				return
			}
			dialer, Err := proxy.FromURL(u, proxy.Direct)
			if Err != nil {
				fmt.Println("创建代理 Dialer 错误:", err)
				return
			}*/
			subdomain := strings.SplitN(addr, ".", 0)
			for _, subdomains := range subdomain {
				switch subdomains {
				case "forum":
					{
						break
					}
				case "github":
					{
						break
					}
				case "mail":
					{
						break
					}
				case "baidu":
					{
						break
					}
				case "bing":
					{
						break
					}
				case "google":
					{
						break
					}
				}
			}
			return proxy, 0
		}
	case 0x04:
		{
			file, err := os.Open("ipv6.txt")
			if err != nil {
				fmt.Println("无法打开文件:", err)
				return proxy, -1
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line[0] == '0' {
					n := len(line)
					var cidr = line[2:n]
					ip, ipNet, _ := net.ParseCIDR(cidr)
					mask := ipNet.Mask.String()
					var ll = ip.String()
					lastIP := make(net.IP, len(ip))
					for i := range ip {
						lastIP[i] = ip[i] | ^mask[i]
					}
					if count > 0 {
						return proxy, count
					}
					var rr = lastIP.String()

					l1 := uint64(ll[0])<<48 | uint64(ll[1])<<32 | uint64(ll[2])<<16 | uint64(ll[3])
					l2 := uint64(ll[4])<<48 | uint64(ll[5])<<32 | uint64(ll[6])<<16 | uint64(ll[7])
					r1 := uint64(rr[0])<<48 | uint64(rr[1])<<32 | uint64(rr[2])<<16 | uint64(rr[3])
					r2 := uint64(rr[4])<<48 | uint64(rr[5])<<32 | uint64(rr[6])<<16 | uint64(rr[7])
					a1 := uint64(addr[0])<<48 | uint64(addr[1])<<32 | uint64(addr[2])<<16 | uint64(addr[3])
					a2 := uint64(addr[4])<<48 | uint64(addr[5])<<32 | uint64(addr[6])<<16 | uint64(addr[7])
					if ((a1 > l1) && (a1 < r1)) || ((a1 == l1) && (a2 > l2)) || ((a1 == r1) && (a2 < r2)) {
						count = 0
					} else {
						count = -1
					}
				} else {
					if count >= 0 {
						n := len(line)
						proxy[count] = line[2:n]
						count++
					}
				}
			}
			return proxy, count
		}
	}
	return proxy, -1
}

func http(array []byte) ([16]string, error, int) {
	search := []byte{0x48, 0x54, 0x54, 0x50, 0x2F, 0x31, 0x2E, 0x31}
	index := bytes.Index(array, search)
	var proxy [16]string
	fmt.Println("http head location:", index)
	fmt.Println("http processing")
	if index == -1 {
		fmt.Println("no http")
		return proxy, errors.New("wrong http"), 0
	}
	search = []byte{0x48, 0x6F, 0x73, 0x74}
	index = bytes.Index(array, search)
	fmt.Println("Host location:", index)
	index += 6
	ss := ""
	for i := index; ; i++ {
		if array[i] == '\n' {
			break
		}
		ss += string(array[i])
	}
	var count = 0
	var Type byte = 0
	if ss[0] >= '0' && ss[0] <= '9' {
		cnt := 0
		for i := 0; i < len(ss); i++ {
			if ss[i] == '.' {
				cnt++
			}
		}
		if cnt == 3 {
			Type = 1
		} else {
			Type = 4
		}
	} else {
		Type = 0x03
	}
	fmt.Println(ss)
	fmt.Println(Type)
	proxy, count = divide(ss, Type)
	return proxy, nil, count
}

func tls(array []byte, n int) ([16]string, error, int) {
	var count = 0
	var proxy [16]string
	fmt.Println("tls processing")
	if array[n] == 0x16 && array[n+1] == 0x03 && array[n+2] == 0x01 {
		ss := ""
		b := array[n+112]
		c := array[n+113]
		d := int(b)*256 + int(c)
		for i := 0; i < d; {
			b = array[n+114+i]
			c = array[n+115+i]
			e := int(b)*256 + int(c)
			if e == 0x00 {
				b = array[n+116+i]
				c = array[n+117+i]
				e = int(b)*256 + int(c)
				for j := 0; j < e; j++ {
					ss += string(array[n+118+i+j])
				}
				var Type byte = 0
				if ss[0] >= '0' && ss[0] <= '9' {
					cnt := 0
					for i := 0; i < len(ss); i++ {
						if ss[i] == '.' {
							cnt++
						}
					}
					if cnt == 3 {
						Type = 1
					} else {
						Type = 4
					}
				} else {
					Type = 0x03
				}
				fmt.Println("tls addr:" + ss)
				proxy, count = divide(ss, Type)
				return proxy, nil, count
			} else {
				b = array[n+116+i]
				c = array[n+117+i]
				e = int(b)*256 + int(c)
				i += e + 4
			}
		}
	}
	fmt.Println("wrong tls")
	return proxy, errors.New("wrong tls"), 0
}

func pid() ([16]string, error, int) {
	var count = 0
	var proxy [16]string
	var inode = 0
	file, err := os.Open("/proc/net/tcp")
	fmt.Println("pid processing")
	if err != nil {
		fmt.Println("wrong tls")
		return proxy, errors.New("wrong tls"), count
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var i = 0
		for ; i < len(line); i++ {
			if line[i] == ':' {
				break
			}
		}
		if i == len(line) {
			continue
		}
		i += 25
		if line[i] == '1' && line[i+1] == 'F' && line[i+2] == '9' && line[i+3] == '0' {
			i += 62
			for {
				if line[i] == ' ' {
					break
				}
				inode += int(line[i] - '0')
			}
		}
	}
	file, err = os.Open("/proc/net/tcp6")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return proxy, errors.New("wrong tls"), count
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var i = 0
		for ; i < len(line); i++ {
			if line[i] == ':' {
				break
			}
		}
		if i == len(line) {
			continue
		}
		i += 73
		if line[i] == '1' && line[i+1] == 'F' && line[i+2] == '9' && line[i+3] == '0' {
			i += 62
			for {
				if line[i] == ' ' {
					break
				}
				inode += int(line[i] - '0')
			}
		}
	}
	dir := "/proc"
	f, Err := os.Open(dir)
	if Err != nil {
		fmt.Println("Failed to open directory:", err)
		return proxy, errors.New("wrong tls"), count
	}
	defer f.Close()
	fileInfos, ERr := f.Readdir(-1)
	if ERr != nil {
		fmt.Println("Failed to read directory:", err)
		return proxy, errors.New("wrong tls"), count
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			Pid, err := strconv.Atoi(fileInfo.Name())
			if err != nil {
				continue
			}
			if Pid < 1000 {
				continue
			}
			filePath := dir + "/" + fileInfo.Name()
			filePath += "/fd"
			ff, Err := os.Open(filePath)
			if Err != nil {
				fmt.Println("Failed to open directory:", err)
				return proxy, errors.New("wrong tls"), count
			}
			fileIn, ERr := ff.Readdir(-1)
			if ERr != nil {
				fmt.Println("Failed to read directory:", err)
				return proxy, errors.New("wrong tls"), count
			}
			for _, fileI := range fileIn {
				if !fileI.IsDir() {
					ss := filePath + "/" + fileI.Name()
					realPath, err := os.Readlink(ss)
					if err != nil {
						fmt.Printf("Failed to read file descriptor %s: %v\n", file.Name(), err)
						continue
					}
					sss := "socks[" + strconv.Itoa(inode) + "]"
					if strings.Contains(realPath, sss) {
						fdDir := fmt.Sprintf("/proc/%d/fd/exe", Pid)
						exeInfo, _ := os.Readlink(fdDir)
						if strings.Contains(exeInfo, "edge") {
							return proxy, errors.New("wrong tls"), count
						}
						if strings.Contains(exeInfo, "mail") {
							return proxy, errors.New("wrong tls"), count
						}
						if strings.Contains(exeInfo, "firefox") {
							return proxy, errors.New("wrong tls"), count
						}
					}
				}
			}
		}
	}
	return proxy, errors.New("wrong tls"), count
}
