setlocal enabledelayedexpansion

SET image_name=grslin/gjjr_kvs
SET ip=172.20.99.161

docker rm -f replica1 replica2 replica3
REM docker rmi %image_name%

docker build -t %image_name% .
docker image prune -f

docker run -d --name=replica1 -p 8080:8080 -e VIEW="%ip%:8080,%ip%:8081,%ip%:8082" -e IP_PORT="%ip%:8080" %image_name%
docker run -d --name=replica2 -p 8081:8080 -e VIEW="%ip%:8080,%ip%:8081,%ip%:8082" -e IP_PORT="%ip%:8081" %image_name%
docker run -d --name=replica3 -p 8082:8080 -e VIEW="%ip%:8080,%ip%:8081,%ip%:8082" -e IP_PORT="%ip%:8082" %image_name%

