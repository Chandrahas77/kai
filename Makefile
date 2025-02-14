build:
	docker-compose build

run:
	docker-compose up --build

stop:
	docker-compose down

logs:
	docker-compose logs -f

psql:
	docker exec -it postgres psql -U user -d vulnerabilities_db