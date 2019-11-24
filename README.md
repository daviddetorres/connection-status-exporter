# Connection status exporter
Exporter of socket connection status for prometheus.

The exporter generates labeled metrics with the status of the socket connection equivalent to execute:
```
nc -z host port
```
## Why a connection status exporter?
Many exporters include per-se connection status metrics for the services that are being monitored. 

Nevertheless, there are some use cases such as IoT devices, smart meters for energy, gas, water or dataloggers that their connection depends on GPRS/3G. 

These IoT devices usually are installed in places where connectivity signal is not strong (electrical rooms, transformation centers, basements, rural areas...) or moving around the territory (like smart vehicle or tracking applications).  

Also in IoT solutions it is common to have only connectivity with the device at certain moments due to they depend on a battery or energy harvesting system. 

Monitoring the connectivity of these devices can help to:
- Detect anomalies and possible failures
- Detect changes in the connectivity profile of the devices (in order to change the GPRS/3G operator or the location of the antenna)
- Predict and detect when a device is available to perform online interventions (updating firmware, download logs, etc.)

## Getting started
To run it:
```
./connection-status-exporter
```
by default the exporter will read the configuration file *config/config.yaml*. Another file can be passed as a flag:
```
./connection-status-exporter --config-file=config/user_config.yaml
```

The metrics are available at http://localhost:9293/metrics. Here is an example: 
```
# HELP connection_status_up Connection status of the socket.
# TYPE connection_status_up gauge
connection_status_up{host="127.0.0.1",name="hostname-http",port="80",protocol="tcp"} 1
connection_status_up{host="localhost",name="hostname-https",port="8080",protocol="udp"} 0
```
The metrics will have the name *connection_status_up* with the labels for each socket and the following possible values:
* 1: Connection OK
* 0: Connection ERROR

## Usage
To configure the sockets that te exporter will check, a yaml configuration file is used. This is a configuration example:
```
sockets:
  - name: hostname-http
    host: 127.0.0.1
    port: 80
    protocol: tcp
  - name: hostname-https 
    host: localhost
    port: 8080
    timeout: 2
```
The fields of the sockets to configure are:
* **name**: A name to be able to filter by this in prometheus
* **host**: Hostname or IP of the socket
* **port**: Port to check
* **protocol**: network parameter of the [Dial function](https://golang.org/pkg/net/#Dial). Known networks are: "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket". If not defined, it will be set to "tcp" by default. 
* **timeout**: Seconds to wait to open the connection until is marked as error. If not specified it is set to 5s.

The following fields will be used as labels in the metric:
* name
* host
* port
* protocol

## Contributing
Please read the [CONTRIBUTING](https://github.com/daviddetorres/connection-status-exporter/blob/master/CONTRIBUTING.md) guidelines.

## Credits
- [David de Torres](https://github.com/daviddetorres) - *Initial work*

## License
This project is published under Apache 2.0, see [LICENSE](https://github.com/daviddetorres/connection-status-exporter/blob/master/LICENSE).