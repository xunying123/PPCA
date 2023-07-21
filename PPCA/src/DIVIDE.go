package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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
					var num uint8 = 0
					var ll uint32 = 0
					var rr uint32 = 0
					i := 2
					for ; i < n; i++ {
						if line[i] == '.' {
							ll = ll*1000 + uint32(num)
							num = 0
						} else if line[i] == '/' {
							ll = ll*1000 + uint32(num)
							num = 0
						} else if line[i] == ' ' {
							break
						} else {
							num = num*10 + line[i] - '0'
						}
					}
					rr = 0b11111111111111111111111111111111
					var temp uint32 = 0b11111111111111111111111111111111
					for j := 0; uint8(j) < num; j++ {
						rr = rr << 1
					}
					ll = ll & rr
					rr = temp ^ rr
					rr = ll | rr
					nn := len(addr)
					var Num uint8 = 0
					var nll uint32 = 0
					for j := 0; j < nn; j++ {
						if addr[j] == '.' {
							nll = nll*1000 + uint32(Num)
							Num = 0
						} else if addr[j] == ':' {
							break
						} else {
							Num = Num*10 + addr[j] - '0'
						}
					}
					if count > 0 {
						return proxy, count
					}
					if (nll > ll) && (nll < rr) {
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
			return proxy, -1
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
					COU := 0
					var num uint8 = 0
					var ll1 uint64 = 0
					var ll2 uint64 = 0
					var rr1 uint64 = 0
					var rr2 uint64 = 0
					i := 0
					for ; i < n; i++ {
						if line[i] == ':' {
							ll2 = ll2*10000 + uint64(num)
							COU++
							if COU == 8 {
								ll1 = ll2
								ll2 = 0
							}
							num = 0
						} else if line[i] == '/' {
							ll2 = ll2*10000 + uint64(num)
							num = 0
						} else if line[i] == ' ' {
							break
						} else {
							num = num*10 + line[i] - '0'
						}
					}
					rr1 = 0b1111111111111111111111111111111111111111111111111111111111111111
					rr2 = 0b1111111111111111111111111111111111111111111111111111111111111111
					var temp uint64 = 0b1111111111111111111111111111111111111111111111111111111111111111
					for j := 0; uint8(j) < num; j++ {
						if j <= 31 {
							rr2 = rr2 << 1
						} else {
							rr1 = rr1 << 1
						}

					}
					ll1 = ll1 & rr1
					ll2 = ll2 & rr2
					rr1 = temp ^ rr1
					rr2 = temp ^ rr2
					rr1 = ll1 | rr1
					rr2 = ll2 | rr2
					nn := len(addr)
					var Num uint8 = 0
					var nll1 uint64 = 0
					var nll2 uint64 = 0
					COU = 0
					for j := 1; j < nn; j++ {
						if addr[j] == ':' {
							nll2 = nll2*10000 + uint64(Num)
							COU++
							if COU == 8 {
								nll1 = nll2
								nll2 = 0
							}
							Num = 0
						} else if addr[j] == ']' {
							break
						} else {
							Num = Num*10 + addr[j] - '0'
						}
					}
					if count > 0 {
						return proxy, count
					}
					if ((nll1 > ll1) && (nll1 < rr1)) || ((nll1 == ll1) && (nll2 > ll2)) || ((nll1 == rr1) && (nll2 < rr2)) {
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
	fmt.Println(index)
	if index == -1 {
		return proxy, errors.New("wrong http"), 0
	}
	search = []byte{0x48, 0x6F, 0x73, 0x74, 0x3A, 0x20, 0x2E}
	index = bytes.Index(array, search)
	index += 6
	ss := ""
	for i := index; ; i++ {
		if array[i] == '\n' {
			break
		}
		ss += string(array[i])
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
	var count = 0
	proxy, count = divide(ss, Type)
	return proxy, nil, count
}

func tls(array []byte, n int) ([16]string, error, int) {
	var count = 0
	var proxy [16]string
	if array[n] == 0x16 && array[n+1] == 0x03 && array[n+2] == 0x01 {
		ss := ""
		a := array[n+110]
		b := array[n+110+int(a)]
		c := array[n+111+int(a)]
		d := int(b)*256 + int(c)
		for i := 0; i < d; {
			b = array[n+112+int(a)+i]
			c = array[n+112+int(a)+i+1]
			e := int(b)*256 + int(c)
			if e == 0x00 {
				b = array[n+112+int(a)+i+2]
				c = array[n+112+int(a)+i+3]
				e = int(b)*256 + int(c)
				for j := 0; j < e; j++ {
					ss += string(array[n+112+int(a)+i+4+j])
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
				proxy, count = divide(ss, Type)
				return proxy, nil, count
			} else {
				b = array[n+112+int(a)+i+2]
				c = array[n+112+int(a)+i+3]
				e = int(b)*256 + int(c)
				i += e + 4
			}
		}
	}
	return proxy, errors.New("wrong tls"), 0
}

func pid() ([16]string, error, int) {
	var count = 0
	var proxy [16]string
	var inode = 0
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		fmt.Println("无法打开文件:", err)
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
