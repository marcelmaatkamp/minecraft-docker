# Minecraft Java Paper Server + Geyser + Floodgate Docker Container
# Author: James A. Chambers - https://jamesachambers.com/minecraft-java-bedrock-server-together-geyser-floodgate/
# GitHub Repository: https://github.com/TheRemote/Legendary-Java-Minecraft-Geyser-Floodgate

# Use Ubuntu rolling version for builder
FROM ubuntu:rolling AS builder

# Update apt
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install qemu-user-static binfmt-support apt-utils -yqq && rm -rf /var/cache/apt/*

# Use Ubuntu rolling version
FROM --platform=linux/amd64 ubuntu:rolling

# Fetch dependencies
RUN apt update && DEBIAN_FRONTEND=noninteractive apt-get install golang openjdk-19-jre-headless systemd-sysv tzdata sudo curl unzip net-tools gawk openssl findutils pigz libcurl4 libc6 libcrypt1 apt-utils libcurl4-openssl-dev ca-certificates binfmt-support nano -yqq && rm -rf /var/cache/apt/*

# Set port environment variable
ENV Port=25565

# Set Bedrock port environment variable
ENV BedrockPort=19132

# Optional maximum memory Minecraft is allowed to use
ENV MaxMemory=

# Optional Paper Minecraft Version override
ENV Version="1.19.4"

# Optional Timezone
ENV TZ="Europe/London"

# Optional folder to ignore during backup operations
ENV NoBackup=""

# Number of rolling backups to keep
ENV BackupCount=10

# Optional switch to skip permissions check
ENV NoPermCheck=""

# Optional switch to tell curl to suppress the progress meter which generates much less noise in the logs
ENV QuietCurl=""

# Optional switch to disable ViaVersion
ENV NoViaVersion=""

ENV BOT_TOKEN=""
ENV APP_ID=""
ENV GUILD_ID=""
ENV OWNER_ID=""
ENV BEDROCK_ADDRESS=""
ENV BEDROCK_PORT=""
ENV JAVA_ADDRESS=""
ENV JAVA_PORT=""
ENV LOGS_CHANNEL=""
ENV START_STOP_TIMEOUT_IN_SECONDS=""
ENV AUTOSTOP_TIMEOUT_IN_MINUTES=""

# IPV4 Ports
EXPOSE 25565/tcp
EXPOSE 19132/tcp
EXPOSE 19132/udp

# Copy scripts to minecraftbe folder and make them executable
RUN mkdir /scripts
RUN mkdir /discordbot
COPY *.sh /scripts/
COPY *.yml /scripts/
COPY server.properties /scripts/
COPY discordbot/ /discordbot/
RUN chmod -R +x /scripts/*.sh
WORKDIR /discordbot
RUN go mod download
RUN go build ./src/bot.go
WORKDIR /

# Set entrypoint to start.sh script
ENTRYPOINT ["/bin/bash", "/scripts/start.sh"]
