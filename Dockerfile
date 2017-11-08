FROM quay.io/gravitational/debian-grande:0.0.1

# bundle the main provisioner program
ADD build/keygen /usr/local/bin/keygen

# By setting this entry point, we expose make target as command
ENTRYPOINT ["/usr/bin/dumb-init", "keygen", "serve"]
