FROM node:22-slim

# Install claude CLI
RUN npm install -g @anthropic-ai/claude-code

# Non-root user for safety
RUN useradd -m -u 1000 agent
USER agent
WORKDIR /workspace

ENTRYPOINT ["claude"]
