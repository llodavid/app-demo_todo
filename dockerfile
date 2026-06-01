FROM alpine:3.21.2
RUN mkdir -p /app/dist
COPY ./dist /app/dist
RUN mkdir -p /app/dist-db
COPY ./dist-db /app/dist-db
WORKDIR /app/dist
LABEL maintainer="https://github.com/RobertTC32/"
ENV OCI_INT_PORT=8080
EXPOSE $OCI_INT_PORT/tcp
RUN chmod +x main.exe
ENTRYPOINT [ "./main.exe" ]
