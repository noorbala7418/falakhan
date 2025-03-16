# Falakhan

Falakhan is an S3 file sender. Receives files from specific paths (you define them in config.yml ) and uploads to S3 object storage. You can define as many independent paths and S3 credentials as you want.

## Build

```bash
docker buildx build -t falakhan:latest . 
```

## Run

```bash
docker run -it --rm -v "$PWD/config.yml:/app/config.yml" falakhan:latest
```

## Config

You can add your paths and s3 credentials in `config.yml` file. Example:

```yml
listen_port: 0.0.0.0:8000

log_level: debug # info, debug, warning, fatal

s3_credentials:
  - name: object1
    access_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    secret_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    endpoint: https://s3.aws.xxx.xxx.xxx.com
    region: us-east-1
    bucket: bucket1

  - name: object2
    access_key: YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
    secret_key: yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy
    endpoint: https://s3.voltr.yyy.yyy.yyy.com
    region: us-east-1
    bucket: bucket2

routes:
  - name: application1 # Route is http://IP:8000/upload/application1
    s3_config: object1

  - name: application2 # Route is http://IP:8000/upload/application1
    s3_config: object2

```

## Test

You can test your application using this commands:

Request:

```bash
# first, we create a 2mb file
fallocate -l 2M falakhan-test

curl -X POST http://localhost:8000/upload/application1 -F file=@falakhan-test
```

Response:

```bash
➜ curl -X POST localhost:8000/upload/application1 -F file=@falakhan-test
{"location":"https://bucket1.s3.aws.xxx.xxx.xxx.com/application1/falakhan-test","status":"done"}
```
