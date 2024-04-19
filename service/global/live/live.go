package live

import (
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"simple-video-net/global"
	"simple-video-net/utils/location"
)

func init() {
	go func() {
		err := Start()
		if err != nil {
			global.Logger.Error("Failed to turn on live service")
		}
	}()
}

func Start() error {
	path := location.GetCurrentAbPath()
	path = path + `\Config\live\`
	cmd := exec.Command("cmd.exe", "/c", "start "+path+"live-go.exe")
	//Get the output object from which you can read the output results
	if stdio, err := cmd.StdoutPipe(); err != nil {
		return err
		log.Fatal(err)
	} else {
		// Guaranteed shutdown of output streams
		defer func(stdio io.ReadCloser) {
			err := stdio.Close()
			if err != nil {

			}
		}(stdio)
		// Run command
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if _, err := ioutil.ReadAll(stdio); err != nil {
			// Read the output
			log.Fatal(err)
			return err
		} else {
			//log.Println(string(opBytes))
		}
	}
	return nil
}
