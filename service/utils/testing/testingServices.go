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

	ipPort := CheckPortsAsLocalHost(liveConfig.IP, []string{"8090", "7001"})
	if len(ipPort) == 0 {
		global.Logger.Info("开启直播")
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
			global.Logger.Warn(ipPort, "Port not open(fail)!")
		} else {
			if conn != nil {
				global.Logger.Info(ipPort, ipPort, "Port is open(success)!")
				err := conn.Close()
				if err != nil {
					return nil
				}
			} else {
				global.Logger.Warn(ipPort, "Port not open(fail)!")
			}
		}
	}
	return nil
}
