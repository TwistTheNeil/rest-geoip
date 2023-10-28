# Build spa
FROM node:18.18.2-alpine3.18 AS spa-builder
WORKDIR /app
RUN npm install -g pnpm
COPY frontend /app
RUN rm -rf /app/node_modules
RUN pnpm install --frozen-lockfile
RUN npx vite build --outDir /app/dist

# Build app
FROM golang:1.21.3-alpine3.17 AS builder
RUN apk add upx
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN rm -rf /app/internal/router/dist
COPY --from=spa-builder /app/dist /app/internal/router/dist
RUN go build -v 
RUN upx /app/rest-geoip

# Main docker image
FROM alpine:3.18.4
COPY --from=builder /app/rest-geoip /usr/bin/
ENV RELEASE_MODE=true
CMD ["/usr/bin/rest-geoip"]
