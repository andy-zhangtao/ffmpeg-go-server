package avfilter
// #cgo pkg-config: libavfilter libavdevice
// #cgo LDFLAGS: libifade.a
// #include "copy.h"
import "C"
import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

/**
IFade 模板
调用`ifade`滤镜生成带有多次淡入效果的视频
*/

type fileInfo struct {
	InFile  string `json:"in_file"`
	OutFile string `json:"out_file"`
}

type Ifade struct{}

func (Ifade) RegisterInfo() (string, func(http.ResponseWriter, *http.Request)) {
	return "ifade", ifade
}

func ifade(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("InValid Parameter %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fi fileInfo

	err = json.Unmarshal(data, &fi)
	if err != nil {
		logrus.Errorf("Unmarshal Parameter Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Infof("in: %s out: %s", fi.InFile, fi.OutFile)

	in_file := fi.InFile
	out_file := fi.OutFile
	ret := C.copy(C.CString(in_file), C.CString(out_file))
	if ret < 0 {
		logrus.Error("IFade Invoke Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
