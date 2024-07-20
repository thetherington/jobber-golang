GO_MODULE := github.com/thetherington/jobber-common

.PHONY: clean
clean:
ifeq ($(OS), Windows_NT)
	if exist "protogen" rd /s /q protogen
	mkdir protogen\go
else
	rm -fR ./protogen 
	mkdir -p ./protogen/go
endif


.PHONY: protoc-go
protoc-go:
	protoc --go_opt=module=${GO_MODULE} --go_out=. \
	--go-grpc_opt=module=${GO_MODULE} --go-grpc_out=. \
	./proto/auth/*.proto \
	./proto/notification/*.proto \
	./proto/users/*.proto \
	./proto/gig/*.proto \
	./proto/chat/*.proto \
	./proto/order/*.proto \
	./proto/review/*.proto



.PHONY: build
build: clean protoc-go