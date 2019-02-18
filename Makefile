build:
		protoc -I. --go_out=plugins=micro:. proto/user/user.proto
		go mod vendor
		git add --all
		git diff-index --quiet HEAD || git commit -a -m 'fix'
		git push origin master

registry:
		docker build -t eu.gcr.io/hprofits/user-service:latest .
		docker push eu.gcr.io/hprofits/user-service:latest

deploy:
	# protoc -I. --go_out=plugins=micro:. proto/user/user.proto
	protoc -I. --gogofaster_out=plugins=micro,\
	Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:\
	. proto/user/user.proto

	easyjson -all $(GOPATH)/src/github.com/gregory-vc/user-service/proto/user/user.pb.go

	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	go mod vendor
	git add --all
	git diff-index --quiet HEAD || git commit -a -m 'fix'
	git push origin master
	docker build -t eu.gcr.io/hprofits/user-service:latest .
	docker push eu.gcr.io/hprofits/user-service:latest
	kubectl replace -f ./deployments/deployment.yml