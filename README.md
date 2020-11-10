
## template

before this... install this

```sh
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-health-probe
```

before this.. get this
```sh
go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

```sh
protoc \
-I . \
-I ${GOPATH}/src \
-I `go list -m -f "{{.Dir}}" github.com/golang/protobuf` \
-I `go list -m -f "{{.Dir}}" google.golang.org/protobuf` \
-I `go list -m -f "{{.Dir}}" github.com/mwitkow/go-proto-validators` \
--go_out . \
--go_opt paths=source_relative \
--go-grpc_out . \
--go-grpc_opt paths=source_relative \
--govalidators_out . \
--grpc-gateway_out . \
--grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt paths=source_relative \
--grpc-gateway_opt generate_unbound_methods=true \
--openapiv2_out . \
--openapiv2_opt logtostderr=true \
--openapiv2_opt generate_unbound_methods=true \
pkg/proto/grpcgwgokit/pb/grpcgwgokit.proto
```

https://github.com/rephus/grpc-gateway-example

