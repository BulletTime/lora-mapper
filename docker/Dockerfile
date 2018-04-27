FROM debian:stretch
RUN mkdir -p /etc/lora-mapper
EXPOSE 8080
VOLUME /etc/lora-mapper
WORKDIR /root/
COPY lora-mapper-linux-arm .
ENTRYPOINT ["./lora-mapper-linux-arm", "--config", "/etc/lora-mapper/lora-mapper.yaml" ,"-v"]