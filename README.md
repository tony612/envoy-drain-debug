```
                                 +---------------+       +---------------+
        |                        |               |       |               |
        |                      ->|    envoy-1    |-------> grpc-server-1 |
        |                   --/  |               |       |               |
        |                --/     +---------------+       +---------------+
        |             --/
+-------v-------+  --/
|               |-/
|  envoy-http   |-\
|               |  --\
+---|-------^---+     --\
    |       |            --\     +---------------+       +---------------+
    |       |               --\  |               |       |               |
    |       |                  ->|    envoy-2    |-------> grpc-server-2 |
    |       |                    |               |       |               |
    |       |                    +---------------+       +---------------+
+---v-------|---+
|               |
|  http-server  |
|               |
+---------------+
```

## Run

Go is needed, and this repo should be in $GOPATH/src/github.com/tony612/envoy-drain-debug

```
cd compose
make build
docker-compose up
wrk -c 20 -t 20 --timeout 3 -d 3600s http://localhost:8080/grpc
docker exec -it compose_envoy-1_1 bash
tail -f /opt/mount-data/envoy-grpc-1.ingress.error.log

# in compose_envoy-1_1
curl -X POST 127.0.0.1:9901/healthcheck/fail # curl needs to be installed
tcpdump -A -s0 -i any -w /opt/mount-data/dump.pcap
```
