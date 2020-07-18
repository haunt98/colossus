package gatewayv1

//go:generate clang-format -i gateway.proto
//go:generate protoc --go_out=plugins=grpc:. gateway.proto
