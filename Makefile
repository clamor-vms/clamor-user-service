.PHONY: setup init build clean pack upload deploy ensure export stop

DIR := ${CURDIR}
DATETIME := $(shell date)

setup:
	docker run -d -p 5000:5000 --restart=always --name registry registry:2

init:
	docker run --rm -v ${DIR}/src:/go/src/skaioskit -w /go/src/skaioskit lushdigital/docker-golang-dep init

build:
	docker run --rm -v ${DIR}/src:/go/src/skaioskit -w /go/src/skaioskit lushdigital/docker-golang-dep ensure
	docker run --rm -v ${DIR}/src:/go/src/skaioskit -w /go/src/skaioskit golang:latest go build -ldflags "-linkmode external -extldflags -static" -o user

clean:
	rm ${DIR}/src/user

pack:
	docker build -f ./Dockerfile -t localhost:5000/skaioskit/user-service .

upload:
	docker push localhost:5000/skaioskit/user-service
	cat k8s/deployment.template.yaml | sed -e 's/{{BUILD_TIME}}/${DATETIME}/g' > deployment.yaml

deploy:
	envsubst < deployment.yaml | kubectl apply -f -

stop:
	kubectl delete deployments,services,pods,pv,pvc,cronjob,job --all
