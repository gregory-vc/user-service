build:
		protoc -I. --go_out=plugins=micro:. proto/user/user.proto
		git add --all
		git diff-index --quiet HEAD || git commit -a -m 'fix'
		git push origin master

run:
		docker run -p 50053:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns user-service
