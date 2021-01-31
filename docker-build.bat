@ECHO OFF
IF "%1"=="NO-CACHE" docker build --no-cache -f Dockerfile.dev --tag atlas-wcc:latest .
IF NOT "%1"=="NO-CACHE" docker build -f Dockerfile.dev --tag atlas-wcc:latest .
