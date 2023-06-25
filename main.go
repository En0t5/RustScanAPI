package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

// to string
func (p *Port) ToString() string {
	return fmt.Sprintf("%s:%s, service: %s, protocol: %s", p.IP, p.Port, p.Service, p.Protocol)
}

func main() {
	r := gin.Default()

	r.POST("/scan/json", func(c *gin.Context) {
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

		ts := time.Now().Unix()
		//randomStr := fmt.Sprintf("%d", rand.Intn(1000))
		filename := fmt.Sprintf("ips_%d.txt", ts)
		fmt.Println(filename)

		go Run(ips_str[:len(ips_str)-1], filename)

		// return status
		c.JSON(http.StatusOK, gin.H{"filename": filename})
	})

	r.POST("/scan/text", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// 处理读取错误...
			return
		}
		bodyString := string(bodyBytes)

		// 处理IP列表
		var ips []IP
		ipstr := strings.Split(bodyString, "\r\n")
		for _, ip := range ipstr {
			if ip == "" {
				continue
			}
			ips = append(ips, IP{IP: ip})
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

		ts := time.Now().Unix()
		//randomStr := fmt.Sprintf("%d", rand.Intn(1000))
		filename := fmt.Sprintf("ips_%d.txt", ts)
		fmt.Println(filename)

		go Run(ips_str[:len(ips_str)-1], filename)

		// return status
		c.JSON(http.StatusOK, gin.H{"filename": filename})
	})

	r.GET("/show/result", func(c *gin.Context) {
		// 获取当前目录下的cache目录路径
		cacheDir := "./cache/"

		// 读取cache目录下的所有文件
		files, err := os.ReadDir(cacheDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 将文件名添加到列表中
		var filenames []string
		for _, file := range files {
			if !file.IsDir() {
				filenames = append(filenames, file.Name())
			}
		}

		// 返回文件名列表
		c.JSON(http.StatusOK, gin.H{"files": filenames})
	})

	r.GET("/download/:filename", func(c *gin.Context) {
		// 获取当前目录下的cache目录路径
		cacheDir := "./cache/"

		// 获取请求参数中的文件名
		filename := c.Param("filename")

		// 规范化文件路径
		filePath := filepath.Clean(filepath.Join(cacheDir, filename))

		// 检查文件路径是否合法
		if !strings.HasPrefix(filePath, "cache") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "非法文件路径"})
			return
		}

		// 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
			return
		}
		defer file.Close()

		// 设置响应头
		stat, err := file.Stat()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件信息"})
			return
		}

		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Length", string(stat.Size()))

		// 发送文件内容给客户端
		http.ServeContent(c.Writer, c.Request, filename, stat.ModTime(), file)
	})

	r.Run(":50500")
}

func Run(targets string, filename string) {
	output := RunRustScan(targets)
	ports := NmapPortDataCleaning(output)
	out_str := ""
	for _, port := range ports {
		fmt.Println(port.ToString())
		out_str += port.ToString() + "\n"
	}
	WriteFile("./cache/"+filename, []byte(out_str))
}

// write file to disk
func WriteFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getIPs(body string) []IP {
	var ips []IP
	for _, addr := range strings.Split(string(body), "\n") {
		ips = append(ips, IP{IP: strings.TrimSpace(addr)})
	}
	return ips
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
