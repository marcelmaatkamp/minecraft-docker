# Minecraft Java Paper Server + Geyser + Floodgate Docker Container
# Author: James A. Chambers - https://jamesachambers.com/minecraft-java-bedrock-server-together-geyser-floodgate/
# GitHub Repository: https://github.com/TheRemote/Legendary-Java-Minecraft-Geyser-Floodgate

# Use Ubuntu rolling version for builder
FROM ubuntu:rolling AS builder

# Update apt
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install qemu-user-static binfmt-support apt-utils -yqq && rm -rf /var/cache/apt/*

# Use Ubuntu rolling version
FROM ubuntu:rolling

# Fetch dependencies
RUN apt update && DEBIAN_FRONTEND=noninteractive apt-get install golang openjdk-19-jre-headless systemd-sysv tzdata sudo curl unzip net-tools gawk openssl findutils pigz libcurl4 libc6 libcrypt1 apt-utils libcurl4-openssl-dev ca-certificates binfmt-support nano -yqq && rm -rf /var/cache/apt/*

# Set port environment variable
ENV Port=25565

# Set Bedrock port environment variable
ENV BedrockPort=19132

# Optional maximum memory Minecraft is allowed to use
ENV MaxMemory=

# OptionalPaper Minecraft Version override
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
ENV QuietCurl="true"

# Optional switch to disable ViaVersion
ENV NoViaVersion=""

# Discord Bot Token
ENV BOT_TOKEN=""

# Discord Application Id
ENV APP_ID=""

# Guild Id (can be blank but it will take a bit for commands to register globally - this is a limitation of discord)
ENV GUILD_ID=""

# The person who 'owns' the bot (the only one who can run the command to create the message that allows for starting/stopping server)
ENV OWNER_ID=""

# The address for accessing the bedrock minecraft server like mc.example.org
ENV BEDROCK_ADDRESS=""

# The port for accessing the bedrock minecraft server 
ENV BEDROCK_PORT="19132"

# The address for accessing the java minecraft server like mc.example.org
ENV JAVA_ADDRESS=""

# The port for accessing the java minecraft server 
ENV JAVA_PORT="25565"

# Minecraft version for discord bot to display
ENV MC_VERSION="1.19.4"

# The discord channel in which to log who starts/stops the server, leave blank for it to be disabled
ENV LOGS_CHANNEL_ID=""

# The timeout between starting/stopping the server
ENV START_STOP_TIMEOUT_IN_SECONDS="30"

# How long it takes for the server to automatically shutdown once empty
ENV AUTOSTOP_TIMEOUT_IN_MINUTES="30"

# Additional fields to add to the embed created by the /minecraft command
# Fieldname:Content,Fieldname:Content
# : seperates fieldname from content
# , seperates fields
ENV ADDITIONAL_MESSAGES_FOR_EMBED=""

# Java default port
EXPOSE 25565/tcp

# Bedrock default port
EXPOSE 19132/tcp
EXPOSE 19132/udp

# Container setup
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
