package aiv1

//go:generate clang-format -i ai.proto
//go:generate protoc --go_out=plugins=grpc:. ai.proto
