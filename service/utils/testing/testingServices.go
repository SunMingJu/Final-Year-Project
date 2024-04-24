package testing

import (
	"simple-video-net/global"
	"simple-video-net/global/config"
	_ "simple-video-net/global/live"
	"simple-video-net/utils/location"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func init() {
	LiveSeverTesting()
	AssetsSeverTesting()
	FFmpegSeverTesting()
}

// LiveSeverTesting
func LiveSeverTesting() {
	//Get the live service port address
	var liveConfig = config.Config.LiveConfig
	CheckPortsAsLocalHost(liveConfig.IP, []string{liveConfig.RTMP, liveConfig.FLV, liveConfig.Api, liveConfig.HLS})
}

//AssetsSeverTesting Runtime static resource storage directory detection
func AssetsSeverTesting() {
	if !location.IsDir("assets") {
		if err := os.MkdirAll("assets", 0775); err != nil {
			log.Fatalf(" initialization assets Directory failed , wrong reason : %s", err.Error())
			global.Logger.Errorf(" initialization assets Directory failed , wrong reason : %s", err.Error())
		}
	}
}

//FFmpegSeverTesting 
func FFmpegSeverTesting() {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		log.Fatalf("ffmpeg Abnormal status, please check and try again.")
		global.Logger.Errorf("ffmpeg Abnormal status, please check and try again.")
	}
}




// CheckPortsAsLocalHost
func CheckPortsAsLocalHost(ip string, Ports []string) []string {
	//unenabled port
	untenablePort := make([]string, 0)
	for _, ipPort := range Ports {
		// detection port
		ipPort = ip + ":" + ipPort
		conn, err := net.DialTimeout("tcp", ipPort, 3*time.Second)
		if err != nil {
			untenablePort = append(untenablePort, ipPort)
			global.Logger.Errorf(ipPort, "Live broadcast service port:%d Unopened(fail)!", ipPort)
		} else {
			if conn != nil {
				
				err := conn.Close()
				if err != nil {
					return nil
				}
			} else {
				global.Logger.Errorf(ipPort, "Live broadcast service port:%d Unopened(fail)!", ipPort)
			}
		}
	}
	return nil
}
