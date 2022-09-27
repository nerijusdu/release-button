run:
	python io/server.py & go run .

dev:
	go run . & (cd web && npm run dev)

build:
	env GOOS=linux GOARCH=arm GOARM=5 go build
	cd web && npm run build
	rm -r dist
	mkdir -p dist/web/dist
	mkdir -p dist/io
	mv release-button dist/
	mv web/dist/* dist/web/dist/
	cp config.yaml dist/
	cp io/*.py dist/io
	cp io/requirements.txt dist/io/
	cp Makefile dist/

prep:
	cd io && pip install -r requirements.txt

prep3:
	cd io && pip3 install -r requirements.txt

run-bin:
	 (cd io && waitress-serve --host 127.0.0.1 server:run) & ./release-button

run-bin3:
	(cd io && python3 server.py) & ./release-button
