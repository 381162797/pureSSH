package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

//Client 连接结构体
type Client struct {
	user     string
	addr     string
	password string
}

func main() {
	client := getConfig("./config.txt")
	session, err := client.conn()
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = start(session)
	if err != nil {
		log.Fatalln(err)
	}
}

func start(session *ssh.Session) (err error) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("linux", 32, 160, modes)
	if err != nil {
		return
	}
	err = session.Shell()
	if err != nil {
		return

	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return
}

func (c *Client) conn() (*ssh.Session, error) {
	conn, err := ssh.Dial("tcp", c.addr, &ssh.ClientConfig{
		User:            c.user,
		Auth:            []ssh.AuthMethod{ssh.Password(c.password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	return conn.NewSession()
}

func getConfig(path string) *Client {
	//读取path文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(path)
			if err != nil {
				fmt.Println("create file err:", err)
				time.Sleep(time.Second * 3)
				os.Exit(0)
			}
			newFile.Close()
			fmt.Printf("请在 %s 文件下写入如下格式配置:\n", path)
			fmt.Println("username 192.168.0.1:22 password")
			time.Sleep(time.Second * 5)
			os.Exit(0)
		}
		fmt.Println("read config err:", err)
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
	str := string(data)
	strS := strings.Split(str, " ")
	if len(strS) != 3 {
		fmt.Println("配置信息有误:", strS)
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
	return &Client{
		user:     strS[0],
		addr:     strS[1],
		password: strS[2],
	}
}
