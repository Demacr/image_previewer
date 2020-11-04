FROM golang AS build
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod tidy # -d -v
RUN go env
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o image_previewer ./cmd/image_previewer/...

FROM alpine
COPY --from=build /build/image_previewer .
CMD ["./image_previewer"]%