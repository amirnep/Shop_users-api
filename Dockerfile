# Sttart from base image 1.12.13:
FROM golang:1.23.0

# Configure the repo url so we can configure our work directory:
ENV REPO_URL=github.com/amirnep/shop

# Setup out $GOPATH
ENV GOPATH=C:\Users\Nematpour\go

ENV APP_PATH=C:\projects\shop\Shop_users-api

# /app/src/github.com/federicoleon/bookstore_items-api/src

# Copy the entire source code from the current directory to $WORKPATH
ENV WORKPATH=$APP_PATH/src
COPY src $WORKPATH
WORKDIR $WORKPATH

RUN go build -o users-api .

# Expose port 8081 to the world:
EXPOSE 8080

CMD ["./users-api"]