FROM quay.io/fedora/fedora:41

# Create a non-root user and group
RUN groupadd -r comfy && useradd -r -g comfy -d /home/comfy comfy

# Create and set up the work directory
RUN mkdir -p /home/comfy && chown -R comfy:comfy /home/comfy

RUN dnf install -y sudo && usermod -aG wheel comfy

RUN echo 'comfy ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN dnf install -y which zsh

# Switch to the non-root user
USER comfy
WORKDIR /home/comfy

RUN touch /home/comfy/.zshrc

# Copy the binary into the container
COPY --chown=comfy:comfy bin/dotcomfy bin/dotcomfy
COPY --chown=comfy:comfy tests/scripts/* tests/scripts/

# TODO: Copy test scenarios that are wrapped as bash scripts

# Default command (optional, replace with your binary execution command if needed)
# CMD ["bin/dotcomfy"]

