FROM alpine:3.8
USER 496
ADD --chown=496:496 ./kubeiql.elf /usr/local/bin/kubeiql
#ADD --chown=496:496 https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VSN:-v1.11.2}/bin/linux/amd64/kubectl /usr/local/bin/kubectl
#RUN chmod 777 /usr/local/bin/kubectl
EXPOSE 8128

CMD ["/usr/local/bin/kubeiql"]


