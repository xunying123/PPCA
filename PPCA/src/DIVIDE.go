package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func divide(addr string, a byte) ([16]string, int) {
	var proxy [16]string
	var count int = -1
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
