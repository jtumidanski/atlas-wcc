if [[ "$1" = "NO-CACHE" ]]
then
   docker build -f Dockerfile.dev --no-cache --tag atlas-wcc:latest .
else
   docker build -f Dockerfile.dev --tag atlas-wcc:latest .
fi
