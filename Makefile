build:
	protoc -I. --go_out=plugins=micro:/Users/tattoor/source/consignment/user-service \
    proto/user/user.proto

run:
	docker run -p 50053:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns user-service
