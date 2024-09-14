package metrics

import (
	"mbgateway/config"
	"net/http"

	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessedMutex sync.Mutex
	opsProcessed      = promauto.NewCounter(prometheus.CounterOpts{
		Name: "all_uses",
		Help: "The total number of active uses",
	})
)

func Incr() {
	opsProcessedMutex.Lock()
	opsProcessed.Inc()
	opsProcessedMutex.Unlock()
}
func StartMonitor() {
	//prometheus.MustRegister(conns)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(config.Cfg.PromPort), nil)
}
