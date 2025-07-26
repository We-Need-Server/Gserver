# Build stage
FROM golang:1.23-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go modules 복사 및 의존성 다운로드
COPY go.mod ./

# 소스코드 복사
COPY . .

# CGO 비활성화하고 정적 바이너리 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o game-server .

# Runtime stage
FROM alpine:latest

# 보안 업데이트 및 ca-certificates 설치
RUN apk --no-cache add ca-certificates tzdata

# 작업 디렉토리 생성
WORKDIR /root/

# 빌드된 바이너리만 복사
COPY --from=builder /app/game-server .

# 포트 노출 (게임서버 포트)
EXPOSE 8080

# 실행 명령
CMD ["./game-server"]