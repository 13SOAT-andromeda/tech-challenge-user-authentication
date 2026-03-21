FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go-gin-lambda/go.mod go-gin-lambda/go.sum ./
RUN go mod download
COPY go-gin-lambda/ ./
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
CMD [ "bootstrap" ]
