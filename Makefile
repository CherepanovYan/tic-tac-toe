
.PHONY: clean run down style rebuild test

all: down run

run:
	@docker-compose up --build -d

down:
	@docker-compose down

clean:
	@docker-compose down --rmi all --volumes --remove-orphans	

rebuild: clean run

style:
	@find . -name "*.go" -type f -exec gofmt -w {} \;	 