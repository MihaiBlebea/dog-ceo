setup: docker 

docker: docker-build docker-run

docker-build:
	docker build -t serbanblebea/dog_ceo:v0.2 .

docker-run:
	docker run --env HTTP_PORT=8081 -p 8081:8081  --name dog_ceo -d serbanblebea/dog_ceo:v0.2 && \
	open http://localhost:8081

docker-destroy:
	docker rm -f dog_ceo