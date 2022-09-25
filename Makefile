run:
	python io/server.py & go run .

dev:
	go run . & (cd web && npm run dev)

build:
	go build
	cd web && npm run build
	mkdir -p dist/web/dist
	mv release-button dist/
	mv web/dist/* dist/web/dist/
	cp config.yaml dist/
