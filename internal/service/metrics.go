package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var processedURLs = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "antifraud_processed_urls_total",
	Help: "Общее количество проверенных URL",
}, []string{"status", "cache_hit"})
