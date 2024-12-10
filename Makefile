buildbuilder: # call first in /EmailSender
	docker build -t "nekkkkitch/docker" -f .\Dockerfile .
stop:
	docker-compose stop \
	&& docker-compose rm 
start: # call second
	docker-compose build --no-cache \
	&& docker-compose up -d
