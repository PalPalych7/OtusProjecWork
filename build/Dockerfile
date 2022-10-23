# Собираем в гошке
FROM golang:1.16.2 as build

ENV BIN_FILE /opt/banner_rotation/banner_rotation-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/mainService/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="PalPalych"
LABEL SERVICE="banner_rotation"
LABEL MAINTAINERS="PalPalych@bk.ru"

ENV BIN_FILE "/opt/banner_rotation/banner_rotation-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/banner_rotation/config.toml
COPY ./configs/config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
