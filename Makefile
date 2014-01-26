all: go js

go:
	cd comm; go build;\
		cd ../routes; go build;\
		cd ../game; go build;\
		cd ..; go build
js:
	browserify public/src/app.js -o public/bundle.js
