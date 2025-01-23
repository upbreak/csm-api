# 1. 베이스 이미지 설정
FROM ubuntu:22.04 AS base-dependencies

# 2. 필수 패키지 설치
RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    make \
    gawk \
    bison \
    python3 \
    gettext \
    unzip \
    libaio1 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 3. Oracle Instant Client 설치
COPY instantclient-basic-linux.x64-21.16.0.0.0dbru.zip /tmp/
RUN unzip /tmp/instantclient-basic-linux.x64-21.16.0.0.0dbru.zip -d /usr/lib/oracle && \
    ln -sfn /usr/lib/oracle/instantclient_21_16 /usr/lib/oracle/instantclient && \
    ln -sfn /usr/lib/oracle/instantclient/libclntsh.so.21.1 /usr/lib/oracle/instantclient/libclntsh.so && \
    rm /tmp/instantclient-basic-linux.x64-21.16.0.0.0dbru.zip

# 4. LD_LIBRARY_PATH 설정
ENV LD_LIBRARY_PATH=/usr/lib/oracle/instantclient:/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH

# 5. 타임존 설정
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 6. 기본 실행 명령 추가
CMD ["tail", "-f", "/dev/null"]
