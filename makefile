CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
PACKAGE=url-checker/cmd/app
retries = 3
timeout = 2000
connectionLimit = 0

build:
	go build -o ${BINDIR}/url-checker ${PACKAGE}

run:
	go run ${PACKAGE} -input=${input} -output=${output} -retries=${retries} -timeout=${timeout} -connectionLimit=${connectionLimit}