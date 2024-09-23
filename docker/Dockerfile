# Use the official code-server image
FROM codercom/code-server:latest

# Set the working directory
WORKDIR /workspace

# Install additional dependencies or tools as needed
RUN sudo apt-get update && \
    sudo apt-get install -y git curl build-essential

# Expose the code-server port
EXPOSE 8080

# Set the entry point to launch code-server
ENTRYPOINT ["code-server", "--bind-addr", "0.0.0.0:8080", "--auth", "none"]
