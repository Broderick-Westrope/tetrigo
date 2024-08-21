FROM scratch
COPY tetrigo /usr/local/bin/tetrigo

# Expose data volume
VOLUME /data

# Expose ports
EXPOSE 53531/tcp

# Set the default command
ENTRYPOINT [ "/usr/local/bin/tetrigo" ]