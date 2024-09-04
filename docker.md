How to use
1. Install docker
2. go to the directory which has the Dockerfile in it (spotify collab backend)
3. To start up the docker compose
   1. Normal:  `docker compose up -d` [-d for detached mode]
   2. Rebuild if any changes since last build: `docker compose up -d --build`
4. Let everything build and run `docker compose ps` to confirm its running (you should see both psql and api)
5. Send requests on 127.0.0.1:8080
6. Check [routes](.\internal\server\routes.go) folder for the routes
7. To bring down the docker compose 
   1. Along with deleting the db: `docker compose down -v`
   2. Without deleting the db: `docker compose down`
8. 