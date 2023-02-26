FROM alpine:3

WORKDIR /app
COPY main /app/main

EXPOSE 8080
CMD ["./main"]