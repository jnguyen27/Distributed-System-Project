setlocal enabledelayedexpansion

SET image_name=grslin/gjjr_kvs
SET net=asgnet
SET view=10.0.0.20:8080,10.0.0.21:8080,10.0.0.22:8080,10.0.0.23:8080,10.0.0.24:8080,10.0.0.25:8080

docker rm -f replica1 replica2 replica3 replica4 replica5 replica6

docker build -t %image_name% .
docker image prune -f

docker run -d --name=replica1 -p 8082:8080 --net=%net% --ip=10.0.0.20 -e VIEW="%view%" -e IP_PORT="10.0.0.20:8080" -e S="2" %image_name%
docker run -d --name=replica2 -p 8083:8080 --net=%net% --ip=10.0.0.21 -e VIEW="%view%" -e IP_PORT="10.0.0.21:8080" -e S="2" %image_name%
docker run -d --name=replica3 -p 8084:8080 --net=%net% --ip=10.0.0.22 -e VIEW="%view%" -e IP_PORT="10.0.0.22:8080" -e S="2" %image_name%
docker run -d --name=replica4 -p 8085:8080 --net=%net% --ip=10.0.0.23 -e VIEW="%view%" -e IP_PORT="10.0.0.23:8080" -e S="2" %image_name%
docker run -d --name=replica5 -p 8086:8080 --net=%net% --ip=10.0.0.24 -e VIEW="%view%" -e IP_PORT="10.0.0.24:8080" -e S="2" %image_name%
docker run -d --name=replica6 -p 8087:8080 --net=%net% --ip=10.0.0.25 -e VIEW="%view%" -e IP_PORT="10.0.0.25:8080" -e S="2" %image_name%

