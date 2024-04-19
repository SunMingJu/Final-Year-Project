package testing

import (
	"net"
	"simple-video-net/global"
	"simple-video-net/global/config"
	"simple-video-net/global/live"
	"time"
)

// LiveSeverTesting
func LiveSeverTesting() {
	//Get the live service port address
	var liveConfig = config.Config.LiveConfig

	ipPort := CheckPortsAsLocalHost(liveConfig.IP, []string{liveConfig.RTMP, liveConfig.FLV, liveConfig.Api, liveConfig.HLS})
	if len(ipPort) == 0 {
		global.Logger.Info("Start live broadcast")
		err := live.Start()
		if err != nil {
			return
		}
	}
}

// CheckPortsAsLocalHost
func CheckPortsAsLocalHost(ip string, Ports []string) []string {
	//unenabled port
	untenablePort := make([]string, 10)
	for _, ipPort := range Ports {
		// detection port
		ipPort = ip + ":" + ipPort
		conn, err := net.DialTimeout("tcp", ipPort, 3*time.Second)
		if err != nil {
			untenablePort = append(untenablePort, ipPort)
			global.Logger.Warnf(ipPort, "Live broadcast service port:%d Unopened(fail)!", ipPort)
		} else {
			if conn != nil {
				
				err := conn.Close()
				if err != nil {
					return nil
				}
			} else {
				global.Logger.Warnf(ipPort, "Live broadcast service port:%d Unopened(fail)!", ipPort)
			}
		}
	}
	return nil
}
