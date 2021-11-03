package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

var check []bool
var stateFilePath = "report"
var speedFilePath = "speedReport"

func isOnline() bool {
	_, err := net.DialTimeout("tcp", "google.com:80", 900*time.Millisecond)
	return err == nil
}

func speedMeasurement() {
	for {
		user, _ := speedtest.FetchUserInfo()
		serverList, _ := speedtest.FetchServerList(user)
		targets, _ := serverList.FindServer([]int{})
		var text string
		for _, s := range targets {
			s.PingTest()
			s.DownloadTest(false)
			s.UploadTest(false)
			text += fmt.Sprintf("%s\nLatency: %s, Download: %f, Upload: %f\n", getTime(), s.Latency, s.DLSpeed, s.ULSpeed)
		}
		f, err := os.OpenFile(speedFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
		f.Close()
		time.Sleep(5 * time.Minute)
	}
}

func appendToFile() {
	for {
		time.Sleep(1 * time.Minute)
		f, err := os.OpenFile(stateFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		var text string
		text = getTime() + "\n"
		for _, bit := range check {
			if bit {
				text += "Y"
			} else {
				text += "N"
			}
		}
		text += "\n"
		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
		f.Close()
		check = check[:0]
	}
}

func getTime() string {
	dt := time.Now()
	return dt.Format("01-02-2006 15:04:05")
}

func main() {
	go speedMeasurement()
	go appendToFile()
	for {
		check = append(check, isOnline())
		time.Sleep(time.Second)
	}
}

