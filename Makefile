.PHONY: test clean cover

GO ?= go

all: test

test:
	${GO} test -v

clean:
	@rm -rf *.out

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
