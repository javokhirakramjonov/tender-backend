# Targets for Makefile

# Command to run PostgreSQL using Docker Compose
run-db:
	docker-compose up -d db

# Command to build and run the Go app using Docker Compose
run:
	docker-compose up -d app

# Command to stop all services
stop:
	docker-compose down

# Command to view logs for all services
logs:
	docker-compose logs -f
