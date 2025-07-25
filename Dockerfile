FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS version
WORKDIR /build
RUN apk add --no-cache git
ADD .git ./.git
RUN git describe --abbrev=0 --tags | tee ./version


FROM node:22-alpine@sha256:5340cbfc2df14331ab021555fdd9f83f072ce811488e705b0e736b11adeec4bb AS builder-web
ADD gui /build/gui
WORKDIR /build/gui
RUN echo "network-timeout 600000" >> .yarnrc
#RUN yarn config set registry https://registry.npm.taobao.org
#RUN yarn config set sass_binary_site https://cdn.npm.taobao.org/dist/node-sass -g
RUN yarn cache clean && yarn && yarn build

FROM golang:1.23-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS builder
ADD service /build/service
WORKDIR /build/service
COPY --from=version /build/version ./
COPY --from=builder-web /build/web server/router/web
RUN export VERSION=$(cat ./version) && CGO_ENABLED=0 go build -ldflags="-X github.com/v2rayA/v2rayA/conf.Version=${VERSION:1} -s -w" -o v2raya .

FROM v2fly/v2fly-core:v5.21.0@sha256:9068d1e6343807b30391adf9a4870f5d3716a92022fa1440cb3dcc364aec4c9f
COPY --from=builder /build/service/v2raya /usr/bin/
# Set default v2ray binary path environment variable
ENV V2RAYA_V2RAY_BIN=/usr/bin/v2ray
RUN wget -O /usr/local/share/v2ray/LoyalsoldierSite.dat https://github.com/2019COVID/dist-v2ray-rules-dat/raw/refs/heads/main/geosite.dat
RUN apk add --no-cache iptables ip6tables tzdata
LABEL org.opencontainers.image.source=https://github.com/v2rayA/v2rayA
EXPOSE 2017
VOLUME /etc/v2raya
ENTRYPOINT ["v2raya"]
