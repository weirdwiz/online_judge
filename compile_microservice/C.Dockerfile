FROM gcc:latest
ARG filename
ARG directory
COPY . directory
WORKDIR directory
RUN gcc -o compiled directory
CMD ["./compiled"]
