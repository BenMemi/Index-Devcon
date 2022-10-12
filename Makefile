
docker_test:
	docker system prune
	docker build -t main .
	docker run --name main main 

docker create:
	docker build -t main . 

docker test exec: 
	# docker kill $(docker ps -q)
	# docker system prune  
	docker build -t tracker .
	docker run --name tracker tracker 
	docker exec -it tracker bash

test:
	docker build -t tracker .
	docker run --name tracker tracker 

clean: 
	docker system prune

Build_and_Push: 
	# Your Project ID here as argument in line! 
	echo "Set PROJECT_ID=YOUR_PROJECT_ID if you have not !!!"
	docker buildx build --platform linux/amd64 -t gcr.io/$(PROJECT_ID)/tracker . 
	docker push gcr.io/$(PROJECT_ID)/tracker


setting up gcloud etc:
	gcloud auth login
	gcloud auth configure-docker

production: 
	# Your Project ID here as argument in line! 
	echo "Set PROJECT_ID=YOUR_PROJECT_ID if you have not !!!"
	# Build for production for arm64 to make it work on GKE
	docker buildx build --platform linux/amd64 -t gcr.io/$(PROJECT_ID)/lens-tracker . 
	docker push gcr.io/$(PROJECT_ID)/my-indexer
	
	



