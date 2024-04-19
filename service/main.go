package main

import (
	_ "simple-video-net/global/database/mysql"
	_ "simple-video-net/global/database/redis"
	"simple-video-net/router"
	"simple-video-net/utils/testing"
)

func main() {
	//Inspection of live streaming services
	testing.LiveSeverTesting()
	//ces
	router.InitRouter()

}
