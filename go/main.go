package main

import (
	"os"
	"strconv"

	"ffmpeg-go-server/web"
	"github.com/sirupsen/logrus"
)

var port = 8000

func init() {
	if os.Getenv("FF_S_PORT") != "" {
		_port, err := strconv.Atoi(os.Getenv("FF_S_PORT"))
		if err == nil {
			port = _port
		}
	}
}

func main() {

	logrus.Infof("FFMPEG-SERVER RUNING IN %d ", port)

	srv := web.FFServer(8000)

	logrus.Fatal(srv.ListenAndServe())
}
