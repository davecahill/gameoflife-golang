all:
	protoc --go_out=. *.proto && go install dave-cahill.com/golserver
