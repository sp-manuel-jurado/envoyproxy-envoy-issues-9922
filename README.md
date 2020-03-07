# envoyproxy-envoy-issues-9922

Related envoy issue:
- https://github.com/envoyproxy/envoy/issues/9922

It runs some e2e tests against next envoy versions running in docker containers:
- envoy:v1.12.0
- envoy:v1.12.1
- envoy:v1.12.2
- envoy:v1.12.3
- envoy:v1.13.0
- envoy:v1.13.1

Ping gRPC service is used:
```
syntax = "proto3";

package sp.rpc;

import "google/protobuf/empty.proto";

service PingService {
    rpc Ping(google.protobuf.Empty) returns (google.protobuf.Empty);
}
```
- Running port: 10005
- Running envoy port: 10000

Notes:
- Take in account that `envoy.filters.http.grpc_http1_reverse_bridge` is enabled for all routes except for `/sp.rpc.PingService/` that is disabled

## Run e2e tests:
Install dependencies:
- [docker](https://docs.docker.com/install/)
- [docker-compose](https://docs.docker.com/compose/install/)
After that, run tests executing:
```
make envoy-e2e
```

## Envoy config used
```
static_resources:
  listeners:
    - name: listener_1
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 10000
      filter_chains:
        - filters:
            - name: envoy.http_connection_manager
              config:
                access_log:
                  - name: envoy.file_access_log
                    config:
                      path: /dev/stdout
                stat_prefix: ingress_http
                route_config:
                  name: ping
                  virtual_hosts:
                    - name: ping
                      domains: ["*"]
                      routes:
                        - match:
                            prefix: "/sp.rpc.PingService/"
                          route:
                            host_rewrite: ping
                            cluster: ping
                            timeout: 2.00s
                          per_filter_config:
                            envoy.filters.http.grpc_http1_reverse_bridge:
                              disabled: true
                http_filters:
                  - name: envoy.filters.http.grpc_http1_reverse_bridge
                    config:
                      content_type: application/grpc+proto
                      withhold_grpc_frames: true
                  - name: envoy.router
  clusters:
    - name: ping
      connect_timeout: 2.00s
      type: strict_dns
      lb_policy: round_robin
      http2_protocol_options: {}
      load_assignment:
        cluster_name: ping
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: ping
                      port_value: 10005
admin:
  access_log_path: /dev/stdout
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 9902
```

## Results
```
test_1           | === RUN   TestServiceServer_Ping
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | --- FAIL: TestServiceServer_Ping (0.07s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0 (0.02s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3 (0.01s)
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0 (0.01s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1 (0.01s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | FAIL
test_1           | FAIL	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping	0.077s
test_1           | ?   	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping/pkg/sp_rpc	[no test files]
test_1           | FAIL
```

### Logs
Tests:
```
test_1           | === RUN   TestServiceServer_Ping
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | --- FAIL: TestServiceServer_Ping (0.07s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3 (0.01s)
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0 (0.01s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1 (0.02s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | FAIL
test_1           | FAIL	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping	0.075s
test_1           | ?   	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping/pkg/sp_rpc	[no test files]
test_1           | FAIL
```
All:
```
make envoy-e2e
docker-compose up
Creating network "envoyproxy-envoy-issues-9922_default" with the default driver
Pulling ping (golang:1.13)...
1.13: Pulling from library/golang
50e431f79093: Pull complete
dd8c6d374ea5: Pull complete
c85513200d84: Pull complete
55769680e827: Pull complete
15357f5e50c4: Pull complete
e2d9b328fba5: Pull complete
f8e0159fc852: Pull complete
Digest: sha256:98d7ddd9e66920d8a95734440efa36b72357d2c8e7f1cf733e7bf2eb526f2d13
Status: Downloaded newer image for golang:1.13
Pulling envoy-v1-12-0 (envoyproxy/envoy:v1.12.0)...
v1.12.0: Pulling from envoyproxy/envoy
e80174c8b43b: Pull complete
d1072db285cc: Pull complete
858453671e67: Pull complete
3d07b1124f98: Pull complete
990fbe3093cb: Pull complete
4da8686dffe3: Pull complete
b816b805e1d5: Pull complete
f3c50b037f2f: Pull complete
149d66218a49: Pull complete
Digest: sha256:80d260d17d39926e5d405713d59e2d98f0aa1b63936f54c9b471a2f39656b6e4
Status: Downloaded newer image for envoyproxy/envoy:v1.12.0
Pulling envoy-v1-12-1 (envoyproxy/envoy:v1.12.1)...
v1.12.1: Pulling from envoyproxy/envoy
e80174c8b43b: Already exists
d1072db285cc: Already exists
858453671e67: Already exists
3d07b1124f98: Already exists
063a772443a9: Pull complete
bc4fb28152f4: Pull complete
0aa06c0f5561: Pull complete
03295e00c780: Pull complete
1e30547d9e51: Pull complete
Digest: sha256:c485247f63de16331fbe56aec8e9d8a83da76ca0fa8dc872691a09b8d99a314f
Status: Downloaded newer image for envoyproxy/envoy:v1.12.1
Pulling envoy-v1-12-2 (envoyproxy/envoy:v1.12.2)...
v1.12.2: Pulling from envoyproxy/envoy
976a760c94fc: Pull complete
c58992f3c37b: Pull complete
0ca0e5e7f12e: Pull complete
f2a274cc00ca: Pull complete
3ed7ba792571: Pull complete
2ddb5874c4a6: Pull complete
90cdac5d331b: Pull complete
95017d1f0c43: Pull complete
a6e9a87907a0: Pull complete
Digest: sha256:b36ee021fc4d285de7861dbaee01e7437ce1d63814ead6ae3e4dfcad4a951b2e
Status: Downloaded newer image for envoyproxy/envoy:v1.12.2
Pulling envoy-v1-12-3 (envoyproxy/envoy:v1.12.3)...
v1.12.3: Pulling from envoyproxy/envoy
fe703b657a32: Pull complete
f9df1fafd224: Pull complete
a645a4b887f9: Pull complete
57db7fe0b522: Pull complete
a72417fef85a: Pull complete
2b449976ec02: Pull complete
d587b7f80926: Pull complete
ab3ff2ee69fc: Pull complete
7cb790b7b011: Pull complete
Digest: sha256:dc4c2ddbfa59fcc0787cb60b5e8fd57628fe5f5111d297abc7c1496100d758b2
Status: Downloaded newer image for envoyproxy/envoy:v1.12.3
Pulling envoy-v1-13-0 (envoyproxy/envoy:v1.13.0)...
v1.13.0: Pulling from envoyproxy/envoy
0a01a72a686c: Pull complete
cc899a5544da: Pull complete
19197c550755: Pull complete
716d454e56b6: Pull complete
27625dfc099a: Pull complete
183dd4baf3b8: Pull complete
60d86dee7350: Pull complete
a518834869f8: Pull complete
2b82b0573013: Pull complete
Digest: sha256:8e63e2c9edbd6cb0d9f6604f1b38e28660dad16327242899652ecab19ec104b9
Status: Downloaded newer image for envoyproxy/envoy:v1.13.0
Pulling envoy-v1-13-1 (envoyproxy/envoy:v1.13.1)...
v1.13.1: Pulling from envoyproxy/envoy
fe703b657a32: Already exists
f9df1fafd224: Already exists
a645a4b887f9: Already exists
57db7fe0b522: Already exists
ac514b0512c1: Pull complete
310b3ffcbbec: Pull complete
26e4b584a10c: Pull complete
4aedd8a49adf: Pull complete
8fd1217f854d: Pull complete
Digest: sha256:11022c45843f3649a8f4ec6a7b5da32c3699a3522b724c42af671229d5bdfacc
Status: Downloaded newer image for envoyproxy/envoy:v1.13.1
Creating envoyproxy-envoy-issues-9922_ping_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-12-3_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-12-1_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-13-1_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-12-0_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-12-2_1 ... done
Creating envoyproxy-envoy-issues-9922_envoy-v1-13-0_1 ... done
Creating envoyproxy-envoy-issues-9922_test_1          ... done
Attaching to envoyproxy-envoy-issues-9922_ping_1, envoyproxy-envoy-issues-9922_envoy-v1-12-3_1, envoyproxy-envoy-issues-9922_envoy-v1-12-2_1, envoyproxy-envoy-issues-9922_envoy-v1-13-1_1, envoyproxy-envoy-issues-9922_envoy-v1-12-0_1, envoyproxy-envoy-issues-9922_envoy-v1-13-0_1, envoyproxy-envoy-issues-9922_envoy-v1-12-1_1, envoyproxy-envoy-issues-9922_test_1
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:249] initializing epoch 0 (hot restart version=11.104)
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:251] statically linked extensions:
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:253]   access_loggers: envoy.file_access_log,envoy.http_grpc_access_log,envoy.tcp_grpc_access_log
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:256]   filters.http: envoy.buffer,envoy.cors,envoy.csrf,envoy.ext_authz,envoy.fault,envoy.filters.http.adaptive_concurrency,envoy.filters.http.dynamic_forward_proxy,envoy.filters.http.grpc_http1_reverse_bridge,envoy.filters.http.grpc_stats,envoy.filters.http.header_to_metadata,envoy.filters.http.jwt_authn,envoy.filters.http.original_src,envoy.filters.http.rbac,envoy.filters.http.tap,envoy.grpc_http1_bridge,envoy.grpc_json_transcoder,envoy.grpc_web,envoy.gzip,envoy.health_check,envoy.http_dynamo_filter,envoy.ip_tagging,envoy.lua,envoy.rate_limit,envoy.router,envoy.squash
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:259]   filters.listener: envoy.listener.http_inspector,envoy.listener.original_dst,envoy.listener.original_src,envoy.listener.proxy_protocol,envoy.listener.tls_inspector
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:262]   filters.network: envoy.client_ssl_auth,envoy.echo,envoy.ext_authz,envoy.filters.network.dubbo_proxy,envoy.filters.network.mysql_proxy,envoy.filters.network.rbac,envoy.filters.network.sni_cluster,envoy.filters.network.thrift_proxy,envoy.filters.network.zookeeper_proxy,envoy.http_connection_manager,envoy.mongo_proxy,envoy.ratelimit,envoy.redis_proxy,envoy.tcp_proxy
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:264]   stat_sinks: envoy.dog_statsd,envoy.metrics_service,envoy.stat_sinks.hystrix,envoy.statsd
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:266]   tracers: envoy.dynamic.ot,envoy.lightstep,envoy.tracers.datadog,envoy.tracers.opencensus,envoy.tracers.xray,envoy.zipkin
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:269]   transport_sockets.downstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:272]   transport_sockets.upstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-0_1  | [2020-03-07 12:34:30.168][1][info][main] [source/server/server.cc:278] buffer implementation: new
envoy-v1-12-0_1  | [2020-03-07 12:34:30.190][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-0_1  | [2020-03-07 12:34:30.194][1][info][main] [source/server/server.cc:344] admin address: 0.0.0.0:9902
envoy-v1-12-0_1  | [2020-03-07 12:34:30.200][1][info][main] [source/server/server.cc:458] runtime: layers:
envoy-v1-12-0_1  |   - name: base
envoy-v1-12-0_1  |     static_layer:
envoy-v1-12-0_1  |       {}
envoy-v1-12-0_1  |   - name: admin
envoy-v1-12-0_1  |     admin_layer:
envoy-v1-12-0_1  |       {}
envoy-v1-12-1_1  | [2020-03-07 12:34:30.356][1][info][main] [source/server/server.cc:249] initializing epoch 0 (hot restart version=11.104)
envoy-v1-12-1_1  | [2020-03-07 12:34:30.357][1][info][main] [source/server/server.cc:251] statically linked extensions:
envoy-v1-12-1_1  | [2020-03-07 12:34:30.358][1][info][main] [source/server/server.cc:253]   access_loggers: envoy.file_access_log,envoy.http_grpc_access_log,envoy.tcp_grpc_access_log
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:249] initializing epoch 0 (hot restart version=11.104)
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:251] statically linked extensions:
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:253]   access_loggers: envoy.file_access_log,envoy.http_grpc_access_log,envoy.tcp_grpc_access_log
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:256]   filters.http: envoy.buffer,envoy.cors,envoy.csrf,envoy.ext_authz,envoy.fault,envoy.filters.http.adaptive_concurrency,envoy.filters.http.dynamic_forward_proxy,envoy.filters.http.grpc_http1_reverse_bridge,envoy.filters.http.grpc_stats,envoy.filters.http.header_to_metadata,envoy.filters.http.jwt_authn,envoy.filters.http.original_src,envoy.filters.http.rbac,envoy.filters.http.tap,envoy.grpc_http1_bridge,envoy.grpc_json_transcoder,envoy.grpc_web,envoy.gzip,envoy.health_check,envoy.http_dynamo_filter,envoy.ip_tagging,envoy.lua,envoy.rate_limit,envoy.router,envoy.squash
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:259]   filters.listener: envoy.listener.http_inspector,envoy.listener.original_dst,envoy.listener.original_src,envoy.listener.proxy_protocol,envoy.listener.tls_inspector
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:262]   filters.network: envoy.client_ssl_auth,envoy.echo,envoy.ext_authz,envoy.filters.network.dubbo_proxy,envoy.filters.network.mysql_proxy,envoy.filters.network.rbac,envoy.filters.network.sni_cluster,envoy.filters.network.thrift_proxy,envoy.filters.network.zookeeper_proxy,envoy.http_connection_manager,envoy.mongo_proxy,envoy.ratelimit,envoy.redis_proxy,envoy.tcp_proxy
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:264]   stat_sinks: envoy.dog_statsd,envoy.metrics_service,envoy.stat_sinks.hystrix,envoy.statsd
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:266]   tracers: envoy.dynamic.ot,envoy.lightstep,envoy.tracers.datadog,envoy.tracers.opencensus,envoy.tracers.xray,envoy.zipkin
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:269]   transport_sockets.downstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-2_1  | [2020-03-07 12:34:29.965][1][info][main] [source/server/server.cc:272]   transport_sockets.upstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-2_1  | [2020-03-07 12:34:29.966][1][info][main] [source/server/server.cc:278] buffer implementation: new
envoy-v1-12-1_1  | [2020-03-07 12:34:30.358][1][info][main] [source/server/server.cc:256]   filters.http: envoy.buffer,envoy.cors,envoy.csrf,envoy.ext_authz,envoy.fault,envoy.filters.http.adaptive_concurrency,envoy.filters.http.dynamic_forward_proxy,envoy.filters.http.grpc_http1_reverse_bridge,envoy.filters.http.grpc_stats,envoy.filters.http.header_to_metadata,envoy.filters.http.jwt_authn,envoy.filters.http.original_src,envoy.filters.http.rbac,envoy.filters.http.tap,envoy.grpc_http1_bridge,envoy.grpc_json_transcoder,envoy.grpc_web,envoy.gzip,envoy.health_check,envoy.http_dynamo_filter,envoy.ip_tagging,envoy.lua,envoy.rate_limit,envoy.router,envoy.squash
envoy-v1-12-1_1  | [2020-03-07 12:34:30.362][1][info][main] [source/server/server.cc:259]   filters.listener: envoy.listener.http_inspector,envoy.listener.original_dst,envoy.listener.original_src,envoy.listener.proxy_protocol,envoy.listener.tls_inspector
envoy-v1-12-1_1  | [2020-03-07 12:34:30.362][1][info][main] [source/server/server.cc:262]   filters.network: envoy.client_ssl_auth,envoy.echo,envoy.ext_authz,envoy.filters.network.dubbo_proxy,envoy.filters.network.mysql_proxy,envoy.filters.network.rbac,envoy.filters.network.sni_cluster,envoy.filters.network.thrift_proxy,envoy.filters.network.zookeeper_proxy,envoy.http_connection_manager,envoy.mongo_proxy,envoy.ratelimit,envoy.redis_proxy,envoy.tcp_proxy
envoy-v1-12-0_1  | [2020-03-07 12:34:30.203][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-12-0_1  | [2020-03-07 12:34:30.204][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:249] initializing epoch 0 (hot restart version=11.104)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:251] statically linked extensions:
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:253]   access_loggers: envoy.file_access_log,envoy.http_grpc_access_log,envoy.tcp_grpc_access_log
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:256]   filters.http: envoy.buffer,envoy.cors,envoy.csrf,envoy.ext_authz,envoy.fault,envoy.filters.http.adaptive_concurrency,envoy.filters.http.dynamic_forward_proxy,envoy.filters.http.grpc_http1_reverse_bridge,envoy.filters.http.grpc_stats,envoy.filters.http.header_to_metadata,envoy.filters.http.jwt_authn,envoy.filters.http.original_src,envoy.filters.http.rbac,envoy.filters.http.tap,envoy.grpc_http1_bridge,envoy.grpc_json_transcoder,envoy.grpc_web,envoy.gzip,envoy.health_check,envoy.http_dynamo_filter,envoy.ip_tagging,envoy.lua,envoy.rate_limit,envoy.router,envoy.squash
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:259]   filters.listener: envoy.listener.http_inspector,envoy.listener.original_dst,envoy.listener.original_src,envoy.listener.proxy_protocol,envoy.listener.tls_inspector
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:262]   filters.network: envoy.client_ssl_auth,envoy.echo,envoy.ext_authz,envoy.filters.network.dubbo_proxy,envoy.filters.network.mysql_proxy,envoy.filters.network.rbac,envoy.filters.network.sni_cluster,envoy.filters.network.thrift_proxy,envoy.filters.network.zookeeper_proxy,envoy.http_connection_manager,envoy.mongo_proxy,envoy.ratelimit,envoy.redis_proxy,envoy.tcp_proxy
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:264]   stat_sinks: envoy.dog_statsd,envoy.metrics_service,envoy.stat_sinks.hystrix,envoy.statsd
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:266]   tracers: envoy.dynamic.ot,envoy.lightstep,envoy.tracers.datadog,envoy.tracers.opencensus,envoy.tracers.xray,envoy.zipkin
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:269]   transport_sockets.downstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:272]   transport_sockets.upstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-3_1  | [2020-03-07 12:34:29.574][1][info][main] [source/server/server.cc:278] buffer implementation: new
envoy-v1-12-2_1  | [2020-03-07 12:34:30.011][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-2_1  | [2020-03-07 12:34:30.015][1][info][main] [source/server/server.cc:344] admin address: 0.0.0.0:9902
envoy-v1-12-2_1  | [2020-03-07 12:34:30.018][1][info][main] [source/server/server.cc:458] runtime: layers:
envoy-v1-12-2_1  |   - name: base
envoy-v1-12-2_1  |     static_layer:
envoy-v1-12-2_1  |       {}
envoy-v1-12-2_1  |   - name: admin
envoy-v1-12-2_1  |     admin_layer:
envoy-v1-12-2_1  |       {}
envoy-v1-12-2_1  | [2020-03-07 12:34:30.018][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-12-2_1  | [2020-03-07 12:34:30.018][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:251] initializing epoch 0 (hot restart version=11.104)
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:253] statically linked extensions:
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.filters: envoy.filters.dubbo.router
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.stats_sinks: envoy.dog_statsd, envoy.metrics_service, envoy.stat_sinks.hystrix, envoy.statsd
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.transport_sockets.upstream: envoy.transport_sockets.alts, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.tap, envoy.transport_sockets.tls, raw_buffer, tls
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.serializers: dubbo.hessian2
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.protocols: dubbo
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.protocols: auto, binary, binary/non-strict, compact, twitter
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.route_matchers: default
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.resource_monitors: envoy.resource_monitors.fixed_heap, envoy.resource_monitors.injected_resource
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.grpc_credentials: envoy.grpc_credentials.aws_iam, envoy.grpc_credentials.default, envoy.grpc_credentials.file_based_metadata
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.filters.http: envoy.buffer, envoy.cors, envoy.csrf, envoy.ext_authz, envoy.fault, envoy.filters.http.adaptive_concurrency, envoy.filters.http.dynamic_forward_proxy, envoy.filters.http.grpc_http1_reverse_bridge, envoy.filters.http.grpc_stats, envoy.filters.http.header_to_metadata, envoy.filters.http.jwt_authn, envoy.filters.http.on_demand, envoy.filters.http.original_src, envoy.filters.http.rbac, envoy.filters.http.tap, envoy.grpc_http1_bridge, envoy.grpc_json_transcoder, envoy.grpc_web, envoy.gzip, envoy.health_check, envoy.http_dynamo_filter, envoy.ip_tagging, envoy.lua, envoy.rate_limit, envoy.router, envoy.squash
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.filters.listener: envoy.listener.http_inspector, envoy.listener.original_dst, envoy.listener.original_src, envoy.listener.proxy_protocol, envoy.listener.tls_inspector
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.filters: envoy.filters.thrift.rate_limit, envoy.filters.thrift.router
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.tracers: envoy.dynamic.ot, envoy.lightstep, envoy.tracers.datadog, envoy.tracers.opencensus, envoy.tracers.xray, envoy.zipkin
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.resolvers: envoy.ip
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.transport_sockets.downstream: envoy.transport_sockets.alts, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.tap, envoy.transport_sockets.tls, raw_buffer, tls
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.retry_host_predicates: envoy.retry_host_predicates.omit_canary_hosts, envoy.retry_host_predicates.previous_hosts
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.transports: auto, framed, header, unframed
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.health_checkers: envoy.health_checkers.redis
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.access_loggers: envoy.file_access_log, envoy.http_grpc_access_log, envoy.tcp_grpc_access_log
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.clusters: envoy.cluster.eds, envoy.cluster.logical_dns, envoy.cluster.original_dst, envoy.cluster.static, envoy.cluster.strict_dns, envoy.clusters.aggregate, envoy.clusters.dynamic_forward_proxy, envoy.clusters.redis
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.filters.udp_listener: envoy.filters.udp_listener.udp_proxy
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.retry_priorities: envoy.retry_priorities.previous_priorities
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.filters.network: envoy.client_ssl_auth, envoy.echo, envoy.ext_authz, envoy.filters.network.dubbo_proxy, envoy.filters.network.kafka_broker, envoy.filters.network.local_ratelimit, envoy.filters.network.mysql_proxy, envoy.filters.network.rbac, envoy.filters.network.sni_cluster, envoy.filters.network.thrift_proxy, envoy.filters.network.zookeeper_proxy, envoy.http_connection_manager, envoy.mongo_proxy, envoy.ratelimit, envoy.redis_proxy, envoy.tcp_proxy
envoy-v1-13-0_1  | [2020-03-07 12:34:30.252][1][info][main] [source/server/server.cc:255]   envoy.udp_listeners: raw_udp_listener
envoy-v1-12-2_1  | [2020-03-07 12:34:30.042][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
envoy-v1-12-2_1  | [2020-03-07 12:34:30.045][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-2_1  | [2020-03-07 12:34:30.046][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-2_1  | [2020-03-07 12:34:30.046][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-3_1  | [2020-03-07 12:34:29.590][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-0_1  | [2020-03-07 12:34:30.267][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener_components.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
ping_1           | go: downloading github.com/golang/protobuf v1.3.4
ping_1           | go: downloading google.golang.org/grpc v1.27.1
envoy-v1-12-1_1  | [2020-03-07 12:34:30.364][1][info][main] [source/server/server.cc:264]   stat_sinks: envoy.dog_statsd,envoy.metrics_service,envoy.stat_sinks.hystrix,envoy.statsd
envoy-v1-12-1_1  | [2020-03-07 12:34:30.364][1][info][main] [source/server/server.cc:266]   tracers: envoy.dynamic.ot,envoy.lightstep,envoy.tracers.datadog,envoy.tracers.opencensus,envoy.tracers.xray,envoy.zipkin
envoy-v1-12-1_1  | [2020-03-07 12:34:30.364][1][info][main] [source/server/server.cc:269]   transport_sockets.downstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-1_1  | [2020-03-07 12:34:30.364][1][info][main] [source/server/server.cc:272]   transport_sockets.upstream: envoy.transport_sockets.alts,envoy.transport_sockets.raw_buffer,envoy.transport_sockets.tap,envoy.transport_sockets.tls,raw_buffer,tls
envoy-v1-12-1_1  | [2020-03-07 12:34:30.364][1][info][main] [source/server/server.cc:278] buffer implementation: new
test_1           | go clean -testcache && go test -v ./...
envoy-v1-12-3_1  | [2020-03-07 12:34:29.597][1][info][main] [source/server/server.cc:344] admin address: 0.0.0.0:9902
envoy-v1-12-0_1  | [2020-03-07 12:34:30.207][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.599][1][info][main] [source/server/server.cc:458] runtime: layers:
envoy-v1-12-3_1  |   - name: base
envoy-v1-12-3_1  |     static_layer:
envoy-v1-12-3_1  |       {}
envoy-v1-12-3_1  |   - name: admin
envoy-v1-12-3_1  |     admin_layer:
envoy-v1-12-3_1  |       {}
envoy-v1-12-3_1  | [2020-03-07 12:34:29.600][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.600][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-12-1_1  | [2020-03-07 12:34:30.383][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-0_1  | [2020-03-07 12:34:30.269][1][info][main] [source/server/server.cc:336] admin address: 0.0.0.0:9902
envoy-v1-12-0_1  | [2020-03-07 12:34:30.210][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
ping_1           | go: extracting github.com/golang/protobuf v1.3.4
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:251] initializing epoch 0 (hot restart version=11.104)
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:253] statically linked extensions:
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.filters.network: envoy.client_ssl_auth, envoy.echo, envoy.ext_authz, envoy.filters.network.dubbo_proxy, envoy.filters.network.kafka_broker, envoy.filters.network.local_ratelimit, envoy.filters.network.mysql_proxy, envoy.filters.network.rbac, envoy.filters.network.sni_cluster, envoy.filters.network.thrift_proxy, envoy.filters.network.zookeeper_proxy, envoy.http_connection_manager, envoy.mongo_proxy, envoy.ratelimit, envoy.redis_proxy, envoy.tcp_proxy
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.protocols: dubbo
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.filters: envoy.filters.dubbo.router
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.filters.listener: envoy.listener.http_inspector, envoy.listener.original_dst, envoy.listener.original_src, envoy.listener.proxy_protocol, envoy.listener.tls_inspector
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.stats_sinks: envoy.dog_statsd, envoy.metrics_service, envoy.stat_sinks.hystrix, envoy.statsd
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.tracers: envoy.dynamic.ot, envoy.lightstep, envoy.tracers.datadog, envoy.tracers.opencensus, envoy.tracers.xray, envoy.zipkin
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.protocols: auto, binary, binary/non-strict, compact, twitter
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.filters.udp_listener: envoy.filters.udp_listener.udp_proxy
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.route_matchers: default
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.retry_host_predicates: envoy.retry_host_predicates.omit_canary_hosts, envoy.retry_host_predicates.previous_hosts
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.resolvers: envoy.ip
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.filters.http: envoy.buffer, envoy.cors, envoy.csrf, envoy.ext_authz, envoy.fault, envoy.filters.http.adaptive_concurrency, envoy.filters.http.dynamic_forward_proxy, envoy.filters.http.grpc_http1_reverse_bridge, envoy.filters.http.grpc_stats, envoy.filters.http.header_to_metadata, envoy.filters.http.jwt_authn, envoy.filters.http.on_demand, envoy.filters.http.original_src, envoy.filters.http.rbac, envoy.filters.http.tap, envoy.grpc_http1_bridge, envoy.grpc_json_transcoder, envoy.grpc_web, envoy.gzip, envoy.health_check, envoy.http_dynamo_filter, envoy.ip_tagging, envoy.lua, envoy.rate_limit, envoy.router, envoy.squash
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.udp_listeners: raw_udp_listener
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.grpc_credentials: envoy.grpc_credentials.aws_iam, envoy.grpc_credentials.default, envoy.grpc_credentials.file_based_metadata
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.transports: auto, framed, header, unframed
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.thrift_proxy.filters: envoy.filters.thrift.rate_limit, envoy.filters.thrift.router
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.health_checkers: envoy.health_checkers.redis
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.access_loggers: envoy.file_access_log, envoy.http_grpc_access_log, envoy.tcp_grpc_access_log
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.resource_monitors: envoy.resource_monitors.fixed_heap, envoy.resource_monitors.injected_resource
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.retry_priorities: envoy.retry_priorities.previous_priorities
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.transport_sockets.upstream: envoy.transport_sockets.alts, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.tap, envoy.transport_sockets.tls, raw_buffer, tls
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.clusters: envoy.cluster.eds, envoy.cluster.logical_dns, envoy.cluster.original_dst, envoy.cluster.static, envoy.cluster.strict_dns, envoy.clusters.aggregate, envoy.clusters.dynamic_forward_proxy, envoy.clusters.redis
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.dubbo_proxy.serializers: dubbo.hessian2
envoy-v1-13-1_1  | [2020-03-07 12:34:30.055][1][info][main] [source/server/server.cc:255]   envoy.transport_sockets.downstream: envoy.transport_sockets.alts, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.tap, envoy.transport_sockets.tls, raw_buffer, tls
envoy-v1-13-0_1  | [2020-03-07 12:34:30.272][1][info][main] [source/server/server.cc:455] runtime: layers:
envoy-v1-13-0_1  |   - name: base
envoy-v1-13-0_1  |     static_layer:
envoy-v1-13-0_1  |       {}
envoy-v1-13-0_1  |   - name: admin
envoy-v1-13-0_1  |     admin_layer:
envoy-v1-13-0_1  |       {}
envoy-v1-12-2_1  | [2020-03-07 12:34:30.064][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-12-2_1  | [2020-03-07 12:34:30.064][1][info][config] [source/server/configuration_impl.cc:117] loading stats sink configuration
envoy-v1-12-2_1  | [2020-03-07 12:34:30.064][1][info][main] [source/server/server.cc:549] starting main dispatch loop
envoy-v1-13-0_1  | [2020-03-07 12:34:30.272][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-13-0_1  | [2020-03-07 12:34:30.274][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-12-1_1  | [2020-03-07 12:34:30.390][1][info][main] [source/server/server.cc:344] admin address: 0.0.0.0:9902
envoy-v1-13-0_1  | [2020-03-07 12:34:30.276][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
envoy-v1-13-1_1  | [2020-03-07 12:34:30.143][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.api.v2.listener.Filter.config' from file listener_components.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
ping_1           | go: extracting google.golang.org/grpc v1.27.1
envoy-v1-12-2_1  | [2020-03-07 12:34:30.069][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:161] cm init: all clusters initialized
envoy-v1-12-2_1  | [2020-03-07 12:34:30.069][1][info][main] [source/server/server.cc:528] all clusters initialized. initializing init manager
envoy-v1-12-2_1  | [2020-03-07 12:34:30.069][1][info][config] [source/server/listener_manager_impl.cc:578] all dependencies initialized. starting workers
envoy-v1-13-0_1  | [2020-03-07 12:34:30.286][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route_components.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-0_1  | [2020-03-07 12:34:30.286][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-0_1  | [2020-03-07 12:34:30.286][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-3_1  | [2020-03-07 12:34:29.602][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
ping_1           | go: downloading golang.org/x/net v0.0.0-20190311183353-d8887717615a
envoy-v1-13-0_1  | [2020-03-07 12:34:30.289][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-12-1_1  | [2020-03-07 12:34:30.392][1][info][main] [source/server/server.cc:458] runtime: layers:
envoy-v1-12-1_1  |   - name: base
envoy-v1-12-1_1  |     static_layer:
envoy-v1-12-1_1  |       {}
envoy-v1-12-1_1  |   - name: admin
envoy-v1-12-1_1  |     admin_layer:
envoy-v1-12-1_1  |       {}
envoy-v1-13-1_1  | [2020-03-07 12:34:30.145][1][info][main] [source/server/server.cc:336] admin address: 0.0.0.0:9902
envoy-v1-12-0_1  | [2020-03-07 12:34:30.211][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
ping_1           | go: downloading google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
envoy-v1-12-3_1  | [2020-03-07 12:34:29.607][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-3_1  | [2020-03-07 12:34:29.607][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-3_1  | [2020-03-07 12:34:29.607][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-1_1  | [2020-03-07 12:34:30.392][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-13-0_1  | [2020-03-07 12:34:30.290][1][info][config] [source/server/configuration_impl.cc:116] loading stats sink configuration
envoy-v1-13-0_1  | [2020-03-07 12:34:30.290][1][info][main] [source/server/server.cc:550] starting main dispatch loop
envoy-v1-13-0_1  | [2020-03-07 12:34:30.291][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:171] cm init: all clusters initialized
envoy-v1-13-0_1  | [2020-03-07 12:34:30.292][1][info][main] [source/server/server.cc:529] all clusters initialized. initializing init manager
envoy-v1-13-0_1  | [2020-03-07 12:34:30.292][1][info][config] [source/server/listener_manager_impl.cc:707] all dependencies initialized. starting workers
envoy-v1-12-3_1  | [2020-03-07 12:34:29.612][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-12-3_1  | [2020-03-07 12:34:29.612][1][info][config] [source/server/configuration_impl.cc:117] loading stats sink configuration
envoy-v1-13-1_1  | [2020-03-07 12:34:30.147][1][info][main] [source/server/server.cc:455] runtime: layers:
envoy-v1-13-1_1  |   - name: base
envoy-v1-13-1_1  |     static_layer:
envoy-v1-13-1_1  |       {}
envoy-v1-13-1_1  |   - name: admin
envoy-v1-13-1_1  |     admin_layer:
envoy-v1-13-1_1  |       {}
envoy-v1-13-1_1  | [2020-03-07 12:34:30.147][1][info][config] [source/server/configuration_impl.cc:62] loading 0 static secret(s)
envoy-v1-13-1_1  | [2020-03-07 12:34:30.147][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-12-0_1  | [2020-03-07 12:34:30.211][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
ping_1           | go: downloading golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
envoy-v1-12-0_1  | [2020-03-07 12:34:30.214][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-13-1_1  | [2020-03-07 12:34:30.166][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.613][1][info][main] [source/server/server.cc:549] starting main dispatch loop
envoy-v1-12-1_1  | [2020-03-07 12:34:30.393][1][info][config] [source/server/configuration_impl.cc:68] loading 1 cluster(s)
envoy-v1-12-0_1  | [2020-03-07 12:34:30.215][1][info][config] [source/server/configuration_impl.cc:117] loading stats sink configuration
ping_1           | go: extracting golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
envoy-v1-13-1_1  | [2020-03-07 12:34:30.188][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route_components.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-1_1  | [2020-03-07 12:34:30.188][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-1_1  | [2020-03-07 12:34:30.188][1][warning][misc] [source/common/protobuf/utility.cc:441] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-1_1  | [2020-03-07 12:34:30.408][1][info][config] [source/server/configuration_impl.cc:72] loading 1 listener(s)
envoy-v1-12-3_1  | [2020-03-07 12:34:29.616][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:161] cm init: all clusters initialized
envoy-v1-12-3_1  | [2020-03-07 12:34:29.616][1][info][main] [source/server/server.cc:528] all clusters initialized. initializing init manager
envoy-v1-12-3_1  | [2020-03-07 12:34:29.616][1][info][config] [source/server/listener_manager_impl.cc:578] all dependencies initialized. starting workers
envoy-v1-13-1_1  | [2020-03-07 12:34:30.197][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-13-1_1  | [2020-03-07 12:34:30.197][1][info][config] [source/server/configuration_impl.cc:116] loading stats sink configuration
ping_1           | go: extracting golang.org/x/net v0.0.0-20190311183353-d8887717615a
envoy-v1-12-0_1  | [2020-03-07 12:34:30.215][1][info][main] [source/server/server.cc:549] starting main dispatch loop
ping_1           | go: downloading golang.org/x/text v0.3.0
envoy-v1-12-1_1  | [2020-03-07 12:34:30.411][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.api.v2.route.Route.per_filter_config' from file route.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-1_1  | [2020-03-07 12:34:30.198][1][info][main] [source/server/server.cc:550] starting main dispatch loop
envoy-v1-12-0_1  | [2020-03-07 12:34:30.219][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:161] cm init: all clusters initialized
envoy-v1-12-0_1  | [2020-03-07 12:34:30.219][1][info][main] [source/server/server.cc:528] all clusters initialized. initializing init manager
envoy-v1-12-0_1  | [2020-03-07 12:34:30.219][1][info][config] [source/server/listener_manager_impl.cc:578] all dependencies initialized. starting workers
envoy-v1-12-1_1  | [2020-03-07 12:34:30.411][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.network.http_connection_manager.v2.HttpFilter.config' from file http_connection_manager.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-13-1_1  | [2020-03-07 12:34:30.202][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:171] cm init: all clusters initialized
envoy-v1-13-1_1  | [2020-03-07 12:34:30.202][1][info][main] [source/server/server.cc:529] all clusters initialized. initializing init manager
envoy-v1-13-1_1  | [2020-03-07 12:34:30.202][1][info][config] [source/server/listener_manager_impl.cc:707] all dependencies initialized. starting workers
envoy-v1-12-1_1  | [2020-03-07 12:34:30.411][1][warning][misc] [source/common/protobuf/utility.cc:282] Using deprecated option 'envoy.config.filter.accesslog.v2.AccessLog.config' from file accesslog.proto. This configuration will be removed from Envoy soon. Please see https://www.envoyproxy.io/docs/envoy/latest/intro/deprecated for details.
envoy-v1-12-1_1  | [2020-03-07 12:34:30.413][1][info][config] [source/server/configuration_impl.cc:97] loading tracing configuration
envoy-v1-12-1_1  | [2020-03-07 12:34:30.413][1][info][config] [source/server/configuration_impl.cc:117] loading stats sink configuration
envoy-v1-12-1_1  | [2020-03-07 12:34:30.413][1][info][main] [source/server/server.cc:549] starting main dispatch loop
envoy-v1-12-1_1  | [2020-03-07 12:34:30.417][1][info][upstream] [source/common/upstream/cluster_manager_impl.cc:161] cm init: all clusters initialized
envoy-v1-12-1_1  | [2020-03-07 12:34:30.418][1][info][main] [source/server/server.cc:528] all clusters initialized. initializing init manager
envoy-v1-12-1_1  | [2020-03-07 12:34:30.418][1][info][config] [source/server/listener_manager_impl.cc:578] all dependencies initialized. starting workers
ping_1           | go: extracting google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
ping_1           | go: extracting golang.org/x/text v0.3.0
test_1           | go: downloading github.com/golang/protobuf v1.3.4
test_1           | go: downloading google.golang.org/grpc v1.27.1
test_1           | go: downloading github.com/stretchr/testify v1.5.1
ping_1           | go: finding github.com/golang/protobuf v1.3.4
ping_1           | go: finding google.golang.org/grpc v1.27.1
ping_1           | go: finding golang.org/x/net v0.0.0-20190311183353-d8887717615a
test_1           | go: extracting github.com/stretchr/testify v1.5.1
test_1           | go: downloading gopkg.in/yaml.v2 v2.2.2
test_1           | go: downloading github.com/pmezard/go-difflib v1.0.0
test_1           | go: downloading github.com/davecgh/go-spew v1.1.0
ping_1           | go: finding google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
test_1           | go: extracting github.com/pmezard/go-difflib v1.0.0
test_1           | go: extracting gopkg.in/yaml.v2 v2.2.2
test_1           | go: extracting github.com/golang/protobuf v1.3.4
test_1           | go: extracting github.com/davecgh/go-spew v1.1.0
ping_1           | go: finding golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
test_1           | go: extracting google.golang.org/grpc v1.27.1
ping_1           | go: finding golang.org/x/text v0.3.0
test_1           | go: downloading golang.org/x/net v0.0.0-20190311183353-d8887717615a
test_1           | go: downloading google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
test_1           | go: downloading golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
test_1           | go: extracting golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
test_1           | go: extracting golang.org/x/net v0.0.0-20190311183353-d8887717615a
test_1           | go: downloading golang.org/x/text v0.3.0
test_1           | go: extracting google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
test_1           | go: extracting golang.org/x/text v0.3.0
test_1           | go: finding github.com/golang/protobuf v1.3.4
test_1           | go: finding google.golang.org/grpc v1.27.1
test_1           | go: finding golang.org/x/net v0.0.0-20190311183353-d8887717615a
test_1           | go: finding golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
test_1           | go: finding google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
test_1           | go: finding golang.org/x/text v0.3.0
test_1           | go: finding github.com/stretchr/testify v1.5.1
test_1           | go: finding github.com/pmezard/go-difflib v1.0.0
test_1           | go: finding github.com/davecgh/go-spew v1.1.0
test_1           | go: finding gopkg.in/yaml.v2 v2.2.2
envoy-v1-12-0_1  | [2020-03-07T12:34:53.303Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 5 3 1 "-" "grpc-go/1.27.1" "22632dda-3ae9-43db-af41-976f7ff09925" "ping" "172.26.0.2:10005"
envoy-v1-12-1_1  | [2020-03-07T12:34:53.320Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 5 5 1 "-" "grpc-go/1.27.1" "91e1906f-c1c4-47ff-a8eb-48598479ea86" "ping" "172.26.0.2:10005"
envoy-v1-12-2_1  | [2020-03-07T12:34:53.331Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 5 3 1 "-" "grpc-go/1.27.1" "a54573ff-ea4d-4d6c-808b-81aa85714450" "ping" "172.26.0.2:10005"
envoy-v1-12-3_1  | [2020-03-07T12:34:53.339Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 5 4 1 "-" "grpc-go/1.27.1" "d4e7dc4b-ab14-4eb0-a0d0-33e40fcbc5bc" "ping" "172.26.0.2:10005"
envoy-v1-13-0_1  | [2020-03-07T12:34:53.348Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 10 2 0 "-" "grpc-go/1.27.1" "499eb56e-9ad6-48d3-87b3-c87458651dc6" "ping" "172.26.0.2:10005"
envoy-v1-13-1_1  | [2020-03-07T12:34:53.354Z] "POST /sp.rpc.PingService/Ping HTTP/2" 200 - 5 10 7 3 "-" "grpc-go/1.27.1" "83aa2e6d-2803-41ca-bb53-8804e04b12e1" "ping" "172.26.0.2:10005"
test_1           | === RUN   TestServiceServer_Ping
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           | === RUN   TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | --- FAIL: TestServiceServer_Ping (0.07s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.0 (0.02s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.1 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.2 (0.01s)
test_1           |     --- PASS: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.12.3 (0.01s)
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0 (0.01s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.0
test_1           |     --- FAIL: TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1 (0.01s)
test_1           |         main_test.go:59:
test_1           |             	Error Trace:	main_test.go:59
test_1           |             	Error:      	Received unexpected error:
test_1           |             	            	rpc error: code = Unknown desc = grpc: client streaming protocol violation: get <nil>, want <EOF>
test_1           |             	Test:       	TestServiceServer_Ping/Testing_envoy_version:_envoy:v1.13.1
test_1           | FAIL
test_1           | FAIL	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping	0.077s
test_1           | ?   	github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping/pkg/sp_rpc	[no test files]
test_1           | FAIL
test_1           | make: *** [Makefile:2: test] Error 1
envoyproxy-envoy-issues-9922_test_1 exited with code 2
```
