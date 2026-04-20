include .env

# @ stops the echo!

run:
	@go run -ldflags "-X main.googleClientSecret=$(GOOGLE_CLIENT_SECRET)" .

build:
	@go build -ldflags "-X main.googleClientSecret=$(GOOGLE_CLIENT_SECRET)" -o pokedex .