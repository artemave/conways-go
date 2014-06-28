all: build

build:
	(cd conway; go build;\
		cd ../synchronized_broadcaster; go build;\
		cd ..; go build)

get_deps:
	for d in dependencies/* ; do \
		(cd $$d && go get); \
	done
