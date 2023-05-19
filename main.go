package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/gin-gonic/gin"
)

type IP struct {
	IP string `json:"ip"`
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

		// c.Status(http.StatusOK)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "ips": ips_str[:len(ips_str)-1], "output": output})
	})

	r.Run(":50500")
}

// cmd run rustscan
func RunRustScan(ips string) string {
	cmd := "rustscan --ulimit 5000 -a " + ips
	command := exec.Command("/bin/sh", "-c", cmd)
	//fmt.Println(cmd)
	output, err := command.Output()
	fmt.Println(string(output))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

// check ip string
func IsIp(ip string) bool {
	ipPattern := regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	if matches := ipPattern.FindStringSubmatch(ip); matches != nil {
		return true
	}
	return false
}
