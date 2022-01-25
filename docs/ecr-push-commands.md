### Retrieve an authentication token and authenticate your Docker client to your registry.
Use the AWS CLI:

```bash
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 825808845683.dkr.ecr.ap-southeast-1.amazonaws.com
```

### Build your Docker image using the following command. For information on building a Docker file from scratch see the instructions here . You can skip this step if your image is already built:
```bash
docker build -t orijinplus-dev .
```

### After the build completes, tag your image so you can push the image to this repository:
```bash
docker tag orijinplus-dev:latest 825808845683.dkr.ecr.ap-southeast-1.amazonaws.com/orijinplus-dev:latest
```

### Run the following command to push this image to your newly created AWS repository:
```bash
docker push 825808845683.dkr.ecr.ap-southeast-1.amazonaws.com/orijinplus-dev:latest
```