package log_test

import (
	"flag"
	alog "github.com/appscode/log"
	"github.com/golang/glog"
	"testing"
)

func init() {
	// flag.Set("logtostderr", "true")
	flag.Set("v", "5")
}

var data = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func BenchmarkAppsCodeLog(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alog.Infoln(data)
	}
}

func BenchmarkGLog(b *testing.B) {
	for n := 0; n < b.N; n++ {
		glog.Infoln(data)
	}
}
