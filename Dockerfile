FROM alpine:3

WORKDIR /app
COPY ./out/main /app/main
RUN chmod +x ./main
RUN mkdir storage
EXPOSE 8080
CMD ["./main"]