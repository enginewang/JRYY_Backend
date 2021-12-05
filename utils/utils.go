package utils

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"os"
	"time"
)

type EmailParam struct {
	EmailAddr string
	Title     string
	Content       string
}

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func sendEmail(ep EmailParam) {
	e := email.NewEmail()
	e.From = "嘉人有约 <tj_jiarenyouyue@163.com>"
	e.To = []string{ep.EmailAddr}
	e.Subject = ep.Title
	e.Text = []byte(ep.Content)
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "tj_jiarenyouyue@163.com", "ZNGJCNQZKEHPMZWH", "smtp.163.com"), &tls.Config{ServerName: "smtp.163.com"})
	if err != nil {
		fmt.Println("报错")
		fmt.Println(err)
	}
}

func SendEmailList(epList []EmailParam) {
	for _, ep := range epList {
		sendEmail(ep)
		fmt.Printf("已发送至%v\n", ep.EmailAddr)
		time.Sleep(5*time.Second)
	}
}