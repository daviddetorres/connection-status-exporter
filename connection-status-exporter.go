// Copyright 2019 David de Torres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

const (

	// Prefix for Prometheus metrics
	namespace = "connection_status"

	// Default timeout in seconds
	defaultTimeout = 5
)

type socketSet struct {
	Sockets []socket `yalm:"sockets"`
}

type socket struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

// SocketSetExporter Exporter of the status of connection
type SocketSetExporter struct {
	statusMetrics *prometheus.GaugeVec
	mutex         sync.Mutex
	sockets       *socketSet
}

// NewSocketSetExporter Creator of SocketSetExporter
func NewSocketSetExporter(configFile string) *SocketSetExporter {

	yalmFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Fatalf("Error while reading configuration file: &%s", err)
	}

	socketSet := socketSet{}
	err = yaml.Unmarshal(yalmFile, &socketSet)
	if err != nil {
		logger.Fatalf("Error parsing config file: %s", err)
	}
	return &SocketSetExporter{
		statusMetrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: namespace,
				Help: "Connection status of the socket.",
			}, []string{"host", "port", "protocol"}),
		sockets: &socketSet,
	}

}

// Init Initialize the statusMetrics
func (exporter *SocketSetExporter) Init() {
	prometheus.MustRegister(exporter.statusMetrics)
}

// Describe Implements interface
func (exporter *SocketSetExporter) Describe(ch chan<- *prometheus.Desc) {

}

var (
	logger     = log.New(os.Stderr, "", log.Lmicroseconds|log.Ltime|log.Lshortfile)
	configFile = flag.String("config-file", "config/config.yaml", "Exporter configuration file.")
)

func main() {
	flag.Parse()

	exporter := NewSocketSetExporter(*configFile)

	exporter.Init()

	logger.Print("Socket exporter initialized.")
}
