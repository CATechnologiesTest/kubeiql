FROM isl-dsdc.ca.com:5000/ca-standard-images/alpine34:latest
USER 496
ADD ./auth /auth
ADD ./*.crt /
EXPOSE 8128

CMD ["/auth"]


