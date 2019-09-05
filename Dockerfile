FROM busybox
COPY test-go-server /
EXPOSE 8080
ENTRYPOINT ["/test-go-server"]
