build:
	go build -o ./bin/golarge .

run:
	go run . $(DIR)

prod:
	go build -ldflags='-w -s' -o ./bin/golarge .

bench:
	go test -bench . -benchmem

debug:
	go run -gcflags='-m -l' . $(DIR)