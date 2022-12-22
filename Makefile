start:
	docker-compose up --build

start-detach:
	docker-compose up --build -d

stop:
	docker-compose down

migration-up:
	$(GOPATH)/bin/goose -dir repository/database/migration postgres "user=postgres password=postgres dbname=julotest sslmode=disable port=5432" up

migration-down:
	$(GOPATH)/bin/goose -dir repository/database/migration postgres "user=postgres password=postgres dbname=julotest sslmode=disable port=5432" down
