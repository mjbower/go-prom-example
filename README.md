# Description
Build a golang application with built in Prometheus metrics, serving json data.
The Build process uses a standard golang image with all build tools included, and generates a static linux binary which runs inside a `scratch` container.  
The advantages of this approach are: - 
- Smaller size
- Quicker runtime
- Uses less resources
- Smaller attack surface

## Build 
Build the container to your local Docker registry 

``` bash
docker build -t goex1 .
```


## Run container
This exposes port 8080 locally,  the routes are logged in the startup message

```bash
docker run -p 8080:8080 goex1
```

```bash
Listening on: -
```

http://localhost:8080/metrics

http://localhost:8080/health

http://localhost:8080/shownodes

