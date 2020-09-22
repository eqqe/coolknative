# coolknative - install Knative and other Kubernetes components

coolknative is CLI to install knative and other components like minio, nats, kourier, fluent, redis, tekton. This is an executable and a Docker image to be used in a Tekton task.

## Get coolknative
(Install kubectl first)

https://github.com/eskersoftware/coolknative/releases

## Usage

Install on any Kubernetes cluster with at least 8 CPUs and 10GB Memory.
```bash
coolknative install tekton
```

```bash
export U=<FILL IN YOUR DOCKER HUB USERNAME>
export P=<FILL IN YOUR DOCKER HUB PASSWORD>
coolknative install cicd -i namespace1 \
    -f namespace1-webservice \
    -f namespace1-readwebservice \
    -f namespace1-asyncwebservice \
    -u $U \
    -p $P
```

```bash
kubectl proxy&
```
Once Proxying you can navigate to

http://localhost:8001/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/

Go to Pipelines > full-install-pipeline > Create +

Select Pipeline Resources.

Click create.


## Enable TLS

To enable HTTPS with TLS, you need a domain name and a wildcard certificate on this domain.

Place in a file tls-crt the wildcard certificate for mydomain.com and in tls-key the private key of this certificate. 
```bash
export D=<FILL IN YOUR DOMAIN NAME example.com>
export I=<FILL IN YOUR PUBLIC IP>
export U=<FILL IN YOUR DOCKER HUB USERNAME>
export P=<FILL IN YOUR DOCKER HUB PASSWORD>
coolknative install cicd -i namespace1 \
    -f namespace1-webservice \
    -f namespace1-readwebservice \
    -f namespace1-asyncwebservice \
    -u $U \
    -p $P \
    --tls-crt-filename tls-crt \
    --tls-key-filename tls-key \
    -w $D \
    --public-ip $I
```

## Pull from a private Git repository

To pull from a private Git repository, you need the address of the ssh server, a private key file ('ssh-privatekey').
```bash
export G=<FILL IN YOUR GIT SSH SERVER ADDRESS ex: ssh.dev.azure.com>
export APPS_GIT=<FILL IN YOUR GIT REPO ADDRESS>
export FILE_RESOURCES_GIT=<FILL IN YOUR GIT RESOURCES REPO ADDRESS>
export U=<FILL IN YOUR DOCKER HUB USERNAME>
export P=<FILL IN YOUR DOCKER HUB PASSWORD>
coolknative install cicd -i namespace1 \
    -f namespace1-webservice \
    -f namespace1-readwebservice \
    -f namespace1-asyncwebservice \
    -u $U \
    -p $P \
    -g $G \
    -y ssh-privatekey \
    --apps-git $APPS_GIT \
    --file-resources-git $FILE_RESOURCES_GIT
```

## Change address of webservices

```
export A=mysubdomain
export U=<FILL IN YOUR DOCKER HUB USERNAME>
export P=<FILL IN YOUR DOCKER HUB PASSWORD>
coolknative install cicd -i namespace1 \
    -f namespace1-webservice \
    -f namespace1-readwebservice \
    -f namespace1-asyncwebservice \
    -u $U \
    -p $P \
    --knative-serving-domain-template {{.Name}}-{{.Namespace}}.{{.Domain}} \
    -a $A
```



### License

MIT

