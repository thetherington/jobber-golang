GATEWAY_BINARY=gatewayService
NOTIFICATION_BINARY=notificationService
AUTH_BINARY=authService
USERS_BINARY=usersService
GIG_BINARY=gigService
CHAT_BINARY=chatService
ORDER_BINARY=orderService
REVIEW_BINARY=reviewService

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

up_build: build_gateway build_notification build_auth build_users build_gig build_chat build_order build_review
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

build: build_gateway build_notification build_auth build_users build_gig build_chat build_order build_review

## build_gateway: builds the gateway binary as a linux executable
build_gateway:
	@echo "Building gateway binary..."
	cd ./1-gateway-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${GATEWAY_BINARY} ./cmd
	@echo "Done!"

## build_notification: builds the notification binary as a linux executable
build_notification:
	@echo "Building notification binary..."
	cd ./2-notification-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${NOTIFICATION_BINARY} ./cmd
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./3-auth-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${AUTH_BINARY} ./cmd
	@echo "Done!"

## build_users: builds the users binary as a linux executable
build_users:
	@echo "Building users binary..."
	cd ./4-users-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${USERS_BINARY} ./cmd
	@echo "Done!"

## build_gig: builds the gigs binary as a linux executable
build_gig:
	@echo "Building gigs binary..."
	cd ./5-gig-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${GIG_BINARY} ./cmd
	@echo "Done!"

## build_chat: builds the chat binary as a linux executable
build_chat:
	@echo "Building chat binary..."
	cd ./6-chat-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${CHAT_BINARY} ./cmd
	@echo "Done!"

## build_order: builds the order binary as a linux executable
build_order:
	@echo "Building order binary..."
	cd ./7-order-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${ORDER_BINARY} ./cmd
	@echo "Done!"

## build_review: builds the review binary as a linux executable
build_review:
	@echo "Building review binary..."
	cd ./8-review-service && env GOOS=linux CGO_ENABLED=0 go build -o build/${REVIEW_BINARY} ./cmd
	@echo "Done!"
