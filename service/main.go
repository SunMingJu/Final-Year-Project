package main

import (
	"simple-video-net/global/database/mysql" 
	_ "simple-video-net/global/database/redis" 
	"simple-video-net/router"
	_ "simple-video-net/utils/socket"  
	_ "simple-video-net/utils/testing"
)

func main() {

	router.InitRouter()

}
