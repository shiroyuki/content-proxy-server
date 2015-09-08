CURR_PATH=`pwd`
GOPATH_DEPS=${CURR_PATH}/lib
GOPATH_LINK=${GOPATH_DEPS}:${CURR_PATH}:
GO_GET=GOPATH=${GOPATH_DEPS} go get
GO_DOC=GOPATH=${GOPATH_LINK} godoc
GO_BUILD=GOPATH=${GOPATH_LINK} go build
GO_MAIN=server.go

build: install_deps compile
	@echo "Done"

install_deps:
	@echo "Installing dependencies..."
	@mkdir -p ${GOPATH_DEPS}
	@${GO_GET} gopkg.in/mgo.v2
	@${GO_GET} github.com/shiroyuki/passerine-go

compile:
	@echo "Compiling..."
	${GO_BUILD} ${GO_MAIN}

run: compile
	@./server

doc_web: install_deps
	${GO_DOC} -http=:6060

clean:
	rm main
