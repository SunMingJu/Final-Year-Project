package live

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"simple-video-net/global"
	"simple-video-net/utils/location"
)

func init() {
	go func() {
		if !global.Config.ProjectConfig.ProjectStates {
			//Test environment to quickly turn on the live service
			if runtime.GOOS == "windows" {
				err := Start()
				if err != nil {
					global.Logger.Error("Failed to start live streaming service")
			}
		}
	}()
}

func Start() error {
	dir, err := location.GetCurrentAbPath()
	if err != nil {
		global.Logger.Errorf("Failed to quickly start livego service in test environment,Failed to obtain current executable file directory")
		return err
	}
	path := dir + `\Config\live\live-go.exe`
	if _, err := os.Stat(path); os.IsNotExist(err) {
		global.Logger.Errorf("Failed to quickly start livego service in test environment document %s does not exist", path)
		return err
	}
	cmd := exec.Command("cmd.exe", "/c", "start "+path)
	if stdio, err := cmd.StdoutPipe(); err != nil {
		return err
		
	} else {
		// Guaranteed shutdown of output streams
		defer func(stdio io.ReadCloser) {
			err := stdio.Close()
			if err != nil {

			}
		}(stdio)
		// Run command
		if err := cmd.Start(); err != nil {
			
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
