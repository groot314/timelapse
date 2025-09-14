package progress

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func ProgressTempSock(totalDuration float64) string {
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	listener, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := listener.Accept()
		if err != nil {
			log.Fatal("listener accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := 0.0

		var (
			cp   float64
			cpMu sync.RWMutex
		)

		go func() {
			ProgressBar(func() float64 {
				cpMu.RLock()
				defer cpMu.RUnlock()
				return cp
			})
		}()
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)

			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cpMu.Lock()
				cp = float64(c) / totalDuration / 1_000_000
				cpMu.Unlock()
			}
			if strings.Contains(data, "progress=end") {
				cp = 1
			}
			if math.Abs(cp-progress) > 0.0001 {
				progress = cp
			}
		}
	}()

	return sockFileName
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func ProbeDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
