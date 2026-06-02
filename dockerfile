FROM alpine:3.21.2
RUN mkdir -p /app/dist
COPY ./dist /app/dist
RUN mkdir -p /app/dist-db
COPY ./dist-db /app/dist-db
WORKDIR /app/dist
EXPOSE 80/tcp
RUN chmod +x main.exe
ENTRYPOINT [ "./main.exe" ]
