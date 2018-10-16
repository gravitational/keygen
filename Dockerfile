FROM quay.io/gravitational/debian-grande:0.0.1

# bundle the main provisioner program
ADD build/keygen /usr/local/bin/keygen
