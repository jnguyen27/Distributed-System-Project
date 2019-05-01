setlocal enabledelayedexpansion

SET image_name=grslin/gjjr_kvs
SET followers=follower1
SET network=asg2net
SET /a port=8083

docker build -t %image_name% .

docker rm -f leader %followers%

docker run --name=leader -p %port%:8080 --net=%network% --ip=10.0.0.20 -d %image_name%
for %%A in (%followers%) do (	
	SET /a port += 1
	docker run --name=%%A -p !port!:8080 --net=%network% -e MAINIP=10.0.0.20:8080 -d %image_name%	
)

	