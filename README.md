# line-googlehome-bot

## usage
```
$ go build
$ GOOGLE_HOME_IP=<googlehome ip> \
  CHANNEL_SECRET=<line channel secret> \
  CHANNEL_TOKEN=<line channel access token> \
  PORT=<app port> \
  VOLUME=<googlehome volume> \
  VOLUME_RESTORE_TIME=<volume restore time> \
  VOICETEXT_API_KEY=<VoiceText WebAPI key> \
  AWS_S3_BUCKET=<AWS s3 bucket name> \
  AWS_ACCESS_KEY_ID=<AWS access key(if necessary)> \
  AWS_SECRET_ACCESS_KEY=<AWS secret access key(if necessary)> \
  ./line-googlehome-bot
```

### environment variables

|name|value|default|
|:--|:--|:--|
|GOOGLE_HOME_IP|google home's IP||
|CHANNEL_SECRET|line channel secret||
|CHANNEL_TOKEN|line channel access token||
|PORT|app port|80|
|VOLUME|google home's volume|your setting volume|
|VOLUME_RESTORE_TIME|wait time to restore volume after playing||
|VOICETEXT_API_KEY|VoiceText WebAPI key||
|AWS_S3_BUCKET|AWS s3 bucket name||
|AWS_ACCESS_KEY_ID|<AWS access key(if necessary)>||
|AWS_SECRET_ACCESS_KEY=|AWS secret access key(if necessary)>||

## Docker image
https://hub.docker.com/r/sminamot/line-googlehome-bot

## k8s
### helm
```
# install
$ helm secrets install -f helm/helm_vars/secrets.yaml --name line-googlehome-bot ./helm/line-googlehome-bot

# update
$ helm secrets upgrade -f helm/helm_vars/secrets.yaml line-googlehome-bot ./helm/line-googlehome-bot

# update secrets
$ helm secrets edit helm/helm_vars/secrets.yaml
```
