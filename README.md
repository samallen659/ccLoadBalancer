# ccLoadBalancer

This is my implementation of the Load Balancer project from [Coding Challenges](https://codingchallenges.fyi/challenges/intro)

## Setup

The configuration of the loadbalancer is handled via a yaml file. Below is an example configuration. The configuration file must be int be in the same directory as the running process

```yaml
#Configuration for ccLoadBalancer
config:
  - name: primary
    listenAddr: localhost:8080
    algorithm: roundrobin 
    endpoints:
      - localhost:8081
      - localhost:8082
  - name: secondary
    listenAddr: localhost:8000
    algorithm: leastconnection
    endpoints:
      - localhost:8001
      - localhost:8002
  - name: tertiary
    listenAddr: localhost:8085
    algorithm: iphash
    endpoints:
      - localhost:8003
```

### Building

To build the load balanacer use the following command

```bash
go build -o ccLoadBalancer ./cmd/main.go
```

### Running

Run the loadbalancer with the following command

```bash
./ccLoadBalancer config_example.yaml
```
