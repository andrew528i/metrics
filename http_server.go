package metrics

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	metrics *Metrics
	cfg *ServerConfig
}

type ServerConfig struct {
	Addr string
}

func NewServer(metrics *Metrics, cfg *ServerConfig) *Server {
	return &Server{metrics: metrics, cfg: cfg}
}

func (s Server) Run() error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		metricType := params["type"]
		metricName := params["name"]

		var (
			metricValue string
			err         error
		)

		switch metricType {
		case "system":
			metricValue, err = s.getSystemMetrics(metricName)
			break

		case "internal":
			metricValue, err = s.getInternalMetrics(metricName)
			break

		default:
			err = fmt.Errorf("unknown metric type: %s", metricType)
			break
		}

		if err != nil {
			w.Write([]byte(err.Error()))

			return
		}

		w.Write([]byte(metricValue))
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{type:[a-zA-Z_]+}/{name:[a-zA-Z_]+}", handler).Methods("GET")

	http.Handle("/", rtr)

	return http.ListenAndServe(s.cfg.Addr, nil)
}

func (s Server) getSystemMetrics(name string) (string, error) {
	value := .0

	switch name {
	case "cpu":
		value = s.metrics.System.CPU()
		break

	case "ram":
		value = s.metrics.System.RAM()
		break

	case "disk":
		value = s.metrics.System.Disk()
		break
	}

	if value == .0 {
		return "", fmt.Errorf("unknown metric name: %s", name)
	}

	return fmt.Sprintf("%.2f", value), nil
}

func (s Server) getInternalMetrics(name string) (string, error) {
	value := ""

	counterValue, err := s.metrics.Internal.Counter.Get(name)
	if err != nil {
		return "", err
	}

	avgCounterValue, err := s.metrics.Internal.AvgCounter.GetAvg(name)
	if err != nil {
		return "", err
	}

	if counterValue != 0 {
		value = fmt.Sprintf("%d", counterValue)
	} else if avgCounterValue != 0 {
		value = fmt.Sprintf("%d", avgCounterValue)
	} else {
		value = fmt.Sprintf("%d", 0)
	}

	return value, nil
}
