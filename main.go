package main

import (
	"fmt"
	"github.com/adithya-sree/logger"
	"github.com/adithya-sree/speed-monitor/dao"
	"github.com/showwin/speedtest-go/speedtest"
	"time"
)

const LogFile = "/var/log/speed-monitor.log"

var out = logger.GetLogger(LogFile, "main")

func main() {
	out.Info("Initializing DAO")
	d, err := dao.NewDao("default", "192.168.0.116")
	if err != nil {
		out.Errorf("Error initializing DAO: [%v]", err)
		return
	}
	for {
		out.Info("Awake")
		out.Info("Collecting network details")

		user, _ := speedtest.FetchUserInfo()
		serverList, _ := speedtest.FetchServerList(user)
		targets, _ := serverList.FindServer([]int{})

		for _, s := range targets {
			out.Info("Finding latency")
			err := s.PingTest()
			if err != nil {
				out.Errorf("Error finding latency: [%v]", err)
			}

			out.Info("Finding DL speed")
			err = s.DownloadTest(false)
			if err != nil {
				out.Errorf("Error finding DL speed: [%v]", err)
			}

			out.Info("Finding UL speed")
			err = s.UploadTest(false)
			if err != nil {
				out.Errorf("Error finding UL speed: [%v]", err)
			}

			now := time.Now()
			result := dao.NetworkSpeedResult{
				Date:     now,
				Upload:   fmt.Sprintf("%f", s.ULSpeed),
				Download: fmt.Sprintf("%f", s.DLSpeed),
				Ping:     fmt.Sprintf("%s", s.Latency),
			}
			out.Infof("Collected network details at [%v] [%v]", now, result)

			out.Info("Inserting result")
			err = d.Insert(result)
			if err != nil {
				out.Errorf("Error inserting result: [%v]", err)
			}
		}

		out.Info("Sleeping")
		time.Sleep(5 * time.Minute)
	}
}
