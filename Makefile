all: build js

build:
	(cd comm; go build;\
		cd ../game; go build;\
		cd ..; go build)
js:
	browserify public/js/app.js -o public/bundle.js

get_deps:
	for d in dependencies/* ; do \
		(cd $$d && go get); \
	done