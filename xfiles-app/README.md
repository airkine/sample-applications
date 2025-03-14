# X-Files Database Web App

## Overview

This is a fun, X-Files-themed web application built with Golang using the Gin framework. The app features a dark, eerie UI with glowing effects and animations. It allows users to explore fictional classified cases and government conspiracies.

## Features

- **Backend:** Golang with Gin framework, using PostgreSQL for data storage
- **Frontend:** HTML, CSS (with animations), and JavaScript
- **Database:** PostgreSQL running inside a Docker container
- **API Endpoints:**
  - `/api/truth` - Returns a message: "The truth is out there..."
  - `/api/cases` - Returns three random classified cases from the database
- **Dark Themed UI:** Inspired by the 90s hacker aesthetic with flickering text and glowing effects
- **Interactive Cases:** Fetch and display conspiracy cases dynamically
- **Dockerized Deployment:** Runs both PostgreSQL and the Go application as containers

## Folder Structure

## Folder Structure

```
/xfiles-app
│-- main.go            # Go backend with PostgreSQL
│-- go.mod             # Go module file
│-- go.sum             # Dependency management
│-- Dockerfile         # Dockerfile for building the Go app
│-- docker-compose.yml # Docker Compose to manage services
│-- k8s                # Kubernetes manifests
│   │-- deployment.yaml  # Deployment for the app and database
│   │-- service.yaml     # Services for the app and database
│   │-- ingress.yaml     # Ingress routing for the app
│-- /templates
│   │-- index.html     # Main frontend file
│-- /static
│   │-- styles.css     # Styling for the eerie X-Files look
│   │-- script.js      # JavaScript for fetching case files
│-- /db
│   │-- init.sql       # SQL file to initialize the database
```

## Installation & Running

1. Clone the repository:
   ```sh
   git clone https://github.com/your-username/xfiles-app.git
   cd xfiles-app
   ```
2. Build and start the containers:
   ```sh
   docker-compose up --build
   ```
3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

## Running PostgreSQL in a Docker Container

We use `docker-compose.yml` to set up **PostgreSQL** and ensure the database is ready before the Go app starts.

### `docker-compose.yml` (Database & App Setup)
```yaml
version: "3.8"

services:
  db:
    image: postgres:15
    container_name: xfiles-db
    restart: always
    environment:
      POSTGRES_DB: xfiles
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U admin -d xfiles"]
      interval: 5s
      timeout: 3s
      retries: 5

  app:
    build: .
    container_name: xfiles-app
    restart: always
    environment:
      DATABASE_URL: postgres://admin:secret@db:5432/xfiles?sslmode=disable
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "8080:8080"

volumes:
  pgdata:
```

## Updated `Dockerfile` for Go & Distroless Deployment

We use a **multi-stage build** to compile the Go app as a static binary and deploy it in a **distroless** image.

```dockerfile
# Use Debian as the base image for building Go
FROM debian:bullseye-slim AS builder

# Install dependencies
RUN apt-get update && apt-get install -y wget tar gcc libc-dev musl-dev

# Set Go version
ENV GO_VERSION=1.23.6

# Download and install Go manually
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set up Go paths
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"

# Create app directory
WORKDIR /app

# Copy Go module files and install dependencies
COPY go.mod go.sum ./
RUN /usr/local/go/bin/go mod tidy

# Copy the entire application source code
COPY . .

# Build a fully static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && chmod +x main

# Use distroless for minimal runtime
FROM gcr.io/distroless/static:latest

WORKDIR /app

# Copy the built binary and templates from the builder
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

## Alternative Local Viewing

If you want to open `index.html` without running the Go server:
- Use a local web server:
  ```sh
  python3 -m http.server 8000
  ```
- Open in a browser: `http://localhost:8000/templates/index.html`

# X-Files Database Web App

## Overview

This is a fun, X-Files-themed web application built with Golang using the Gin framework. The app features a dark, eerie UI with glowing effects and animations. It allows users to explore fictional classified cases and government conspiracies.

## Features

- **Backend:** Golang with Gin framework, using PostgreSQL for data storage
- **Frontend:** HTML, CSS (with animations), and JavaScript
- **Database:** PostgreSQL running inside a Docker container
- **API Endpoints:**
  - `/api/truth` - Returns a message: "The truth is out there..."
  - `/api/cases` - Returns three random classified cases from the database
- **Dark Themed UI:** Inspired by the 90s hacker aesthetic with flickering text and glowing effects
- **Interactive Cases:** Fetch and display conspiracy cases dynamically
- **Kubernetes Deployment:** Uses Kubernetes for orchestration and NGINX Ingress for routing



## Installation & Running

1. Clone the repository:
   ```sh
   git clone https://github.com/your-username/xfiles-app.git
   cd xfiles-app
   ```
2. Deploy to Kubernetes:
   ```sh
   kubectl apply -f k8s/
   ```
3. Open your browser and visit:
   ```
   http://your-ingress-domain-or-ip
   ```

## Kubernetes Deployment

### Deployment Manifest (`k8s/deployment.yaml`)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: xfiles-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: xfiles-app
  template:
    metadata:
      labels:
        app: xfiles-app
    spec:
      containers:
      - name: xfiles-app
        image: your-dockerhub-username/xfiles-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          value: "postgres://admin:secret@xfiles-db:5432/xfiles?sslmode=disable"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: xfiles-db
spec:
  replicas: 1
  selector:
    matchLabels:
      app: xfiles-db
  template:
    metadata:
      labels:
        app: xfiles-db
    spec:
      containers:
      - name: xfiles-db
        image: postgres:15
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: "xfiles"
        - name: POSTGRES_USER
          value: "admin"
        - name: POSTGRES_PASSWORD
          value: "secret"
        volumeMounts:
        - mountPath: /var/lib/postgresql/data
          name: db-data
      volumes:
      - name: db-data
        emptyDir: {}
```

### Service Manifest (`k8s/service.yaml`)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: xfiles-app
spec:
  selector:
    app: xfiles-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: xfiles-db
spec:
  selector:
    app: xfiles-db
  ports:
  - protocol: TCP
    port: 5432
    targetPort: 5432
  type: ClusterIP
```

### Ingress Manifest (`k8s/ingress.yaml`)
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: xfiles-ingress
spec:
  rules:
  - host: xfiles.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: xfiles-app
            port:
              number: 80
```

## Deploying to Kubernetes

1. **Apply the manifests**:
   ```sh
   kubectl apply -f k8s/
   ```
2. **Check the deployment**:
   ```sh
   kubectl get pods
   ```
3. **Access the app** using your ingress domain:
   ```sh
   http://xfiles.local
   ```

## Future Improvements

- Add user authentication for submitting new cases
- Implement a pagination feature for exploring case archives
- Enhance animations for a more immersive experience

👽 **The truth is out there!** 🚀



