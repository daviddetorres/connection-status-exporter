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
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

const (

	// Prefix for Prometheus metrics
	namespace = "connection_status_up"

	// Default values for optional parameters of socket
	defaultTimeout         = 5
	defaultProtocol string = "tcp"

	// Constant values
	connectionOk          = 1
	connectionErr         = 0
	metricsPublishingPort = ":8888"
)

var (
	logger     = log.New(os.Stderr, "", log.Lmicroseconds|log.Ltime|log.Lshortfile)
	configFile = flag.String("config-file", "config/config.yaml", "Exporter configuration file.")
	addr       = flag.String("listen-address", metricsPublishingPort, "The address to listen on for HTTP requests.")
)

type socketSet struct {
	Sockets []socket `yalm:"sockets"`
}

type socket struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Timeout  int    `yalm:"timeout"`
}

// SocketSetExporter Exporter of the status of connection
type SocketSetExporter struct {
	socketStatusMetrics *prometheus.GaugeVec
	mutex               sync.Mutex
	sockets             *socketSet
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

	err = socketSet.check()
	if err != nil {
		logger.Fatalf("Error in the configuration of the sockets: %s", err)
	}

	return &SocketSetExporter{
		socketStatusMetrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: namespace,
				Help: "Connection status of the socket.",
			}, []string{"name", "host", "port", "protocol"}),
		sockets: &socketSet,
	}
}

// Describe Implements interface
func (exporter *SocketSetExporter) Describe(prometheusChannel chan<- *prometheus.Desc) {
	exporter.socketStatusMetrics.Describe(prometheusChannel)
	return
}

// Collect Implements interface
func (exporter *SocketSetExporter) Collect(prometheusChannel chan<- prometheus.Metric) {
	exporter.mutex.Lock()
	defer exporter.mutex.Unlock()
	exporter.sockets.collect(exporter.socketStatusMetrics)
	exporter.socketStatusMetrics.Collect(prometheusChannel)
	return
}

// Calls the method collect of each socket in the socketSet
func (thisSocketSet *socketSet) collect(prometheusGaugeVector *prometheus.GaugeVec) {
	for _, currentSocket := range thisSocketSet.Sockets {
		currentSocket.collect(prometheusGaugeVector)
	}
	return
}

// Checks the status of the connection of a socket and updates it in the Metric
func (thisSocket *socket) collect(prometheusGaugeVector *prometheus.GaugeVec) {
	connectionStatus := connectionOk
	connectionAdress := thisSocket.Host + ":" + strconv.Itoa(thisSocket.Port)
	connectionTimeout := time.Duration(thisSocket.Timeout * 1000000000)

	// Create a connection to test the socket
	connection, err := net.DialTimeout(
		thisSocket.Protocol,
		connectionAdress,
		connectionTimeout)

	// If the socket cannot be opened, set the status to error
	if err != nil {
		connectionStatus = connectionErr
	}

	// Updated the status of the socket in the metric
	prometheusGaugeVector.WithLabelValues(
		thisSocket.Name,
		thisSocket.Host,
		strconv.Itoa(thisSocket.Port),
		thisSocket.Protocol).Set(float64(connectionStatus))

	// If the socket was open correctly, close it
	if connectionStatus == connectionOk {
		err = connection.Close()
		if err != nil {
			logger.Printf("Error closing the socket")
		}
	}
	return
}

// check the sanity of the sockets in the set
func (thisSocketSet *socketSet) check() error {
	for index := range thisSocketSet.Sockets {
		err := thisSocketSet.Sockets[index].check()
		if err != nil {
			return (err)
		}
	}
	return (nil)
}

// Check the sanity of the socket and fills the default values
func (thisSocket *socket) check() error {
	if thisSocket.Name == "" {
		return (errors.New("All sockets must have the field name completed"))
	}
	if thisSocket.Name == "" {
		return (errors.New("All sockets must have the fiels host completed"))
	}
	if thisSocket.Port == 0 {
		return (errors.New("All sockets must have the field port completed"))
	}
	if thisSocket.Protocol == "" {
		thisSocket.Protocol = defaultProtocol
	}
	// Check if the protocol is among the valid ones
	if IsValidProtocol(thisSocket.Protocol) == false {
		return (errors.New("The protocol of the socket is not a valid one"))
	}
	if thisSocket.Timeout == 0 {
		thisSocket.Timeout = defaultTimeout
	}
	return (nil)
}

// IsValidProtocol Check if a string is among the valid protocols
func IsValidProtocol(protocol string) bool {
	switch protocol {
	case
		"tcp",
		"tcp4",
		"tcp6",
		"udp",
		"udp4",
		"udp6",
		"ip",
		"ip4",
		"ip6",
		"unix",
		"unixgram",
		"unixpacket":
		return true
	}
	return false
}

func main() {
	flag.Parse()
	exporter := NewSocketSetExporter(*configFile)
	prometheus.MustRegister(exporter)
	logger.Print("Socket exporter initialized.")

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
