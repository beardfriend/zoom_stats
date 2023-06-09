# 베이스 이미지로 Node.js 사용
FROM node:latest

# 작업 디렉토리 생성 및 설정
WORKDIR /app

# 앱 종속성 설치
COPY package.json .
COPY package-lock.json .
RUN npm install

# 소스 코드 복사
COPY . .

# React 앱 실행
CMD ["npm", "start"]
