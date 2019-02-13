build:
		protoc -I. --go_out=plugins=micro:. proto/user/user.proto
		go mod vendor
		git add --all
		git diff-index --quiet HEAD || git commit -a -m 'fix'
		git push origin master

registry:
		docker build -t eu.gcr.io/hprofits/user-service:latest .
		gcloud docker -- push eu.gcr.io/hprofits/user-service:latest

deploy:
	protoc -I. --go_out=plugins=micro:. proto/user/user.proto
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	go mod vendor
	git add --all
	git diff-index --quiet HEAD || git commit -a -m 'fix'
	git push origin master
	docker build -t eu.gcr.io/hprofits/user-service:latest .
	gcloud docker -- push eu.gcr.io/hprofits/user-service:latest
	kubectl replace -f ./deployments/deployment.yml