build:
	echo "Building lambda binaries"
	env GOOS=linux GOARCH=arm64 go build -o build/apache-kafka-producer/bootstrap main.go

zip:
	zip -j build/apache-kafka-producer.zip build/apache-kafka-producer/bootstrap