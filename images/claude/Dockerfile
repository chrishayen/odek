FROM node:22-slim

# Install claude CLI
RUN npm install -g @anthropic-ai/claude-code

# node:22-slim ships with a non-root user `node` at uid 1000
USER node
WORKDIR /workspace

ENTRYPOINT ["claude"]
