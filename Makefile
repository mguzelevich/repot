.PHONY: run test clean

exec_root = $(shell pwd)
out = repot
out_path = .build
cache_path = .cache

GOOS = GOOS=linux
GOARCH = GOARCH=amd64

CGO_ENABLED = CGO_ENABLED=1
build_tags = -tags sqlite3
ld_flags = -ldflags '-extldflags -static'

dummy:
	@echo "repot"

help:
	@echo "repot"

fmt:
	gofmt -w cmd fs git helpers repot_test workerpool *.go

clean:
	mkdir -p ${out_path} || rm -rf ${out_path}/*
	mkdir -p ${cache_path}

docker-build-builder:
	docker build -t mguzelevich/golang -f Dockerfile-golang .

build-fast: clean fmt
	go build -o ${out_path}/${out} ./cmd/repot/

build-docker: clean fmt # docker-build-builder
	time docker run -it --rm \
		-v "$(exec_root)/../repot:/repot" \
		-v "$(exec_root):/src" \
		-v "$(exec_root)/${out_path}:/build" \
		-v "$(exec_root)/${cache_path}:/root/go" \
		mguzelevich/golang \
		sh -c "${CGO_ENABLED} ${GOOS} ${GOARCH} go build -v ${build_tags} ${ld_flags} -o /build/${out} ."
