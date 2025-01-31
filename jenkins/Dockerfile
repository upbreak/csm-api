FROM jenkins/jenkins:lts

# 컨테이너 내부에서 root로 필수 패키지 설치
USER root

# 필수 패키지 설치 및 Node.js 설치
RUN apt-get update && apt-get install -y \
    sudo \
    curl \
    sshpass \
    docker.io \
    && rm -rf /var/lib/apt/lists/* \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@10.9.0

# Docker 그룹 및 jenkins 사용자 설정
RUN getent group docker || groupadd -g 999 docker \
    && usermod -aG docker jenkins \
    && chown -R jenkins:jenkins /var/jenkins_home

# 필요한 디렉토리 생성 및 권한 설정
RUN mkdir -p /tmp/csm-front /tmp/csm-api /tmp/base-image \
    && chmod -R 755 /tmp/csm-front /tmp/csm-api /tmp/base-image \
    && chown -R jenkins:jenkins /tmp/csm-front /tmp/csm-api /tmp/base-image

# jenkins 사용자에게 sudo 권한 부여
RUN echo "jenkins ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/jenkins \
    && chmod 440 /etc/sudoers.d/jenkins

# 작업 디렉토리 설정
WORKDIR /var/jenkins_home

# jenkins 사용자로 설정
USER jenkins
