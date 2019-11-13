package cache

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"time"
	"vgo/service"
)

/*
发送验证邮件给注册邮件
缓存于 redis ，设置过期时间
验证邮件由加密密码的 hashcode 组成
hashcode(password) -- id
*/
var from string = os.Getenv("EMAIL")

const hostname = "smtp.qq.com"
const port = 465

var password = os.Getenv("EMAILPASSWORD")

func SendRegisterMail(register *service.UserRegisterService) *service.Response {

	subject := "VGo验证邮件"

	// 内容体
	md5Password := md5.Sum([]byte(register.Password))
	hashPassword := hex.EncodeToString(md5Password[:])
	url := os.Getenv("CHECK_URL") + hashPassword
	body := `<h1>请复制并访问以下网址完成注册验证！</h1><a>` + url + `</a>`

	err := sendEmail(register.Email, subject, body)

	if err != nil {
		return &service.Response{
			Error: err.Error(),
			Code:  53001,
			Msg:   "发送注册邮件失败",
		}
	} else {
		// 发送验证邮件成功
		// 将注册数据先进行缓存, 5分钟后过时
		reJs, _ := json.Marshal(&register)
		RedisClient.Set(hashPassword, reJs, 300*time.Second)
		return &service.Response{Code: 20000, Msg: "发送注册邮件成功"}
	}
}

// 发送验证码给邮箱
func SendCode(email string) *service.Response {
	subject := "VGo验证邮件"
	// 内容体
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	body := `<h1>验证码` + code + `</h1><a>`

	err := sendEmail(email, subject, body)
	if err != nil {
		return &service.Response{
			Error: err.Error(),
			Code:  53002,
			Msg:   "发送验证码邮件失败",
		}
	} else {
		// 发送验证邮件成功
		// 将邮箱进行缓存, 5分钟后过时
		RedisClient.Set(code, email, 300*time.Second)
		return &service.Response{Code: 20000, Msg: "发送验证码邮件成功！"}
	}
}

//return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendRegisterMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// 发送邮件
func sendEmail(toEmail, subject, body string) (err error) {
	auth := smtp.PlainAuth("", from, password, hostname)
	to := []string{toEmail}

	headers := make(map[string]string)
	headers["From"] = "YuGo" + "<" + from + ">"
	headers["To"] = to[0]
	headers["content-Type"] = "text/html; charset=UTF-8"
	headers["Subject"] = subject

	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	err = sendMailUsingTLS(fmt.Sprintf("%s:%d", hostname, port),
		auth, from, to, []byte(msg))
	return err
}
