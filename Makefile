
run:
	go run main.go

build-mac:
	go build -o gitwebhook main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gitwebhook-linux main.go

deploy: build-linux deploy-linux clean-linux

deploy-linux:build-linux
	scp ./gitwebhook-linux root@192.168.1.100:/usr/bin/gitwebhook
	scp -r ./tmpl/*.tmpl root@192.168.1.100:/data/deployment/githook-tmpl/

deploy-mac:build-mac
	cp ./gitwebhook /usr/local/bin/gitwebhook

clean-linux:
	rm -fr ./gitwebhook-linux

clean-mac:
	rm -fr ./gitwebhook

clean:clean-linux clean-mac

