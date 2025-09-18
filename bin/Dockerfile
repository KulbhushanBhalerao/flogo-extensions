# You can replace or modify this Dockerfile as per your requirements
# DO NOT CHANGE or REMOVE 'COPY ./flogoapp /app/flogoapp'
FROM alpine:latest
WORKDIR /app
COPY ./flogoapp /app/flogoapp
RUN chmod +x flogoapp
ENTRYPOINT ["/app/flogoapp"]
