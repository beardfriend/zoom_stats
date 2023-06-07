# Base 이미지 설정
FROM golang:latest

# 작업 디렉토리 설정
WORKDIR /app

# Go 종속성 복사
COPY go.mod go.sum ./

# 종속성 다운로드
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN go build -o main cmd/app/main.go

# 실행 파일 설정
CMD ["./main"]
