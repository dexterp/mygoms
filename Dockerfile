# kick:render
FROM alpine:3.14

RUN apt --no-cache update; apt --no-cache upgrade; apk --no-cache add bash musl-nscd

COPY dist/${PROJECT_NAME}_linux_amd64/${PROJECT_NAME}server /usr/local/bin/
COPY .docker/start.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/${PROJECT_NAME}server
RUN chmod 755 /usr/local/bin/start.sh

ENTRYPOINT ["/usr/local/bin/start.sh"]
