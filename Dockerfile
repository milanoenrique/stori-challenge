# Build Stage
FROM golang:alpine3.20 AS builder

WORKDIR /

# Instalar las dependencias necesarias como make y gcc
RUN apk add --no-cache make gcc g++ musl-dev

# Copiar los archivos de mod y suma para instalar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar todos los archivos del proyecto al contenedor
COPY . .  

# Ejecutar el build usando el Makefile
RUN make build  

# Final Stage
FROM alpine:3.20

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /build/job /PAYMENT-PROCESS 

COPY conf.json /conf.json

# Definir el comando de inicio
CMD ["/PAYMENT-PROCESS"]
