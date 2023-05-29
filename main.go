package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"regexp"
)

type IP struct {
	IP string `json:"ip"`
}

type Port struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Service  string `json:"service"`
	Protocol string `json:"protocol"`
}

func main() {
	r := gin.Default()

	r.POST("/scan", func(c *gin.Context) {
		var ips []IP
		if err := c.BindJSON(&ips); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 处理IP列表
		ips_str := ""
		for _, ip := range ips {
			if !IsIp(ip.IP) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ip is error"})
				return
			}
			ips_str += ip.IP + ","
		}
		fmt.Println(ips_str[:len(ips_str)-1])
		output := RunRustScan(ips_str[:len(ips_str)-1])

		ports := NmapPortDataCleaning(output)

		// c.Status(http.StatusOK)
		c.JSON(http.StatusOK, ports)
	})

	r.Run(":50500")
}

// cmd run rustscan
func RunRustScan(ips string) []byte {
	cmd := "rustscan --ulimit 5000 -a " + ips
	command := exec.Command("/bin/sh", "-c", cmd)
	//fmt.Println(cmd)
	output, err := command.Output()
	fmt.Println(string(output))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return output
}

// check ip string
func IsIp(ip string) bool {
	ipPattern := regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	if matches := ipPattern.FindStringSubmatch(ip); matches != nil {
		return true
	}
	return false
}

// nmap port data clean
func NmapPortDataCleaning(output []byte) []Port {
	//result, _ := os.Open(path)
	// 创建正则表达式匹配模式
	ipPattern := regexp.MustCompile(`Nmap scan report for.*?(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))`)
	portPattern := regexp.MustCompile(`(\d+)/(tcp|udp)\s+open\s+([\w\-\.\+]+)\s*`)
	// 逐行读取文件并匹配ip和端口信息
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var ip string
	var ports []Port
	for scanner.Scan() {
		line := scanner.Text()
		port := Port{}
		// 匹配ip地址
		if matches := ipPattern.FindStringSubmatch(line); matches != nil {
			if matches[1] != "" {
				ip = matches[1]
			}
		}
		port.IP = ip
		// 匹配端口信息
		if matches := portPattern.FindStringSubmatch(line); matches != nil {
			portvalue := matches[1]
			protocol := matches[2]
			service := matches[3]
			port.Port = portvalue
			port.Protocol = protocol
			port.Service = service
			// add port to ports
			ports = append(ports, port)
		}
	}
	return ports
}
