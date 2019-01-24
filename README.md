# multicast-http-proxy
originally use in sync data between elasticsearch clusters; can be use in many scenarios;

### Usage
1. go build
2. modify conf/proxy.conf
3. ./multicast-http-proxy -conf=./conf/proxy.conf &
4. send your es http request to this proxy
