.PHONY: protos

protos:
	protoc -I protos/ protos/btd.proto --go_out=plugins=grpc:protos