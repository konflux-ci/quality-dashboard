FROM golang:1.19 AS builder

WORKDIR /github.com/redhat-appstudio-qe/qe-dashboard-backend
COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.8-860

WORKDIR /root/
COPY --from=builder /github.com/redhat-appstudio-qe/qe-dashboard-backend/bin/server-runtime ./
CMD ["/root/server-runtime"]
