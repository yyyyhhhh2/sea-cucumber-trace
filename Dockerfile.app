FROM scratch

WORKDIR /app

COPY backend/build/server-linux-amd64 /app/server
COPY frontend/dist /app/public

ENV STATIC_DIR=/app/public

EXPOSE 8080

ENTRYPOINT ["/app/server"]
