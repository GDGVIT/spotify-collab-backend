How to use
1. Install docker
2. go to the directory which has the Dockerfile in it (spotify collab backend)
3. run `docker compose up -d` [-d for detached mode]
4. Let everything build and run `docker compose ps` to confirm its running (you should see both psql and api)
5. Send requests on 127.0.0.1:8080
6. Check [routes](.\internal\server\routes.go) folder for the routes