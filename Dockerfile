FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/bootstrap /var/runtime/bootstrap
COPY --from=public.ecr.aws/datadog/lambda-extension:latest /opt/extensions/datadog-agent /opt/extensions/datadog-agent
CMD [ "bootstrap" ]
