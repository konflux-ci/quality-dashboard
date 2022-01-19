FROM golang:1.17 AS builder

WORKDIR /github.com/flacatus/qe-dashboard-backend
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM registry.access.redhat.com/ubi8-minimal:8.4-212

WORKDIR /root/
COPY --from=builder /github.com/flacatus/qe-dashboard-backend/bin/qe-dashboard-backend ./
CMD ["/root/qe-dashboard-backend"]
