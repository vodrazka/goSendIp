package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/bjdean/gonetcheck"
	"github.com/go-gomail/gomail"
	"net/mail"
	"os"
)

var topic string
var targetMail string
var password string
var smtpServer string
var username string
var smtpPort int

func initFlags() {
	flag.StringVar(&topic, "topic", "", "Email subject additional info.")
	flag.StringVar(&targetMail, "target", "", "Target email address.")
	flag.StringVar(&password, "pass", "password", "Source emial password.")
	flag.StringVar(&smtpServer, "smtp", "smtp.mailgun.org", "Address of smtp server.")
	flag.StringVar(&username, "user", "rasp@sandbox85e310238ee2448888290d9f25179241.mailgun.org", "Source email username.")
	flag.IntVar(&smtpPort, "port", 587, "smtp port.")
}

func checkFlags(){
	_, err := mail.ParseAddress(targetMail)
	if err != nil{
		fmt.Printf("Target email incorrect: '%v'", targetMail)
		os.Exit(1)
	}
}

func localAddresses() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return nil
	}
	results := make([]string, 0)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					results = append(results, fmt.Sprintf("%v (MAC:%v): %v<br>", i.Name, i.HardwareAddr, ipnet.IP.String()))
				}
			}

		}
	}
	return results
}

func sendEmail(subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "raspff@gmail.com")
	m.SetHeader("To", targetMail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpServer, smtpPort, username, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Printf("Email sent to %v\n", targetMail)
}

func checkInternetConnection() (bool, []error) {
	return gonetcheck.CheckInternetAccess(
		time.Duration(60 * time.Second),
		[]string{"http://google.com"},
		[]string{})
}

func getSSID() string {
	cmd, err := exec.Command("/sbin/iwgetid").Output()
	if err != nil {
		log.Println(err)
		return "/sbin/iwgetid did not handle"
	}
	return string(cmd)
}

func main() {
	initFlags()
	flag.Parse()
	checkFlags();
	hasInternet, _ := checkInternetConnection()
	if hasInternet {
		addresses := localAddresses()
		if addresses != nil {
			fmt.Printf("%v\n", addresses)
			sendEmail(fmt.Sprintf("%v - %v (%v)", "RaspberryPi - IP", topic, time.Now().Format("Monday 15:04:05 2-01-2006")), strings.Trim(fmt.Sprintf("%v",
				addresses), "[]") + fmt.Sprintf("Wireless network: %v", getSSID()))
		}
	}
}
