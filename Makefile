build:
		protoc -I. --go_out=plugins=micro:. proto/user/user.proto
		go mod vendor
		git add --all
		git diff-index --quiet HEAD || git commit -a -m 'fix'
		git push origin master