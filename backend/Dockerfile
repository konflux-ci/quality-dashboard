FROM golang:1.19 AS builder

WORKDIR /github.com/konflux-ci-qe/qe-dashboard-backend
COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build

FROM registry.access.redhat.com/ubi8-minimal:8.6-854

WORKDIR /root/
COPY --from=builder /github.com/konflux-ci-qe/qe-dashboard-backend/bin/server-runtime ./
CMD ["/root/server-runtime"]
