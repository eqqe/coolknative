FROM alpine:latest

WORKDIR /app
ADD coolknative /app/coolknative
RUN apk --no-cache add ca-certificates && apk add curl && curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl

CMD ["/bin/sh"]
