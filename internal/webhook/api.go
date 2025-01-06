package webhook

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/provider/webhook/api"
)

func Run(provider provider.Provider, serveaddr string, serveaddrMetrics string, readTimeout time.Duration, writeTimeout time.Duration) {
	timeout := time.Duration(1 * time.Minute)
	ch := make(chan struct{}, 1)
	log.Printf("Starting webhook server on %s", serveaddr)
	go api.StartHTTPApi(provider, ch, timeout, timeout, serveaddr)
	<-ch
	log.Printf("Starting metrics server on %s", serveaddrMetrics)

	RunMetrics(serveaddrMetrics)
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

func RunMetrics(serveaddr string) {
	h := http.NewServeMux()
	s := &http.Server{
		Addr:    serveaddr,
		Handler: h,
	}
	h.Handle("/healthz", http.HandlerFunc(okHandler))
	h.Handle("/readyz", http.HandlerFunc(okHandler))
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
