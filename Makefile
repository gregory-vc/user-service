build:
		protoc -I. --go_out=plugins=micro:. proto/user/user.proto
		go mod vendor
		git add --all
		git diff-index --quiet HEAD || git commit -a -m 'fix'
		git push origin master

registry:
		docker build -t eu.gcr.io/my-project-tattoor/user-service:latest .
		gcloud docker -- push eu.gcr.io/my-project-tattoor/user-service:latest

deploy:
	protoc -I. --go_out=plugins=micro:. proto/user/user.proto
	go mod vendor
	git add --all
	git diff-index --quiet HEAD || git commit -a -m 'fix'
	git push origin master
	docker build -t eu.gcr.io/my-project-tattoor/user-service:latest .
	gcloud docker -- push eu.gcr.io/my-project-tattoor/user-service:latest
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl replace -f ./deployments/deployment.yml