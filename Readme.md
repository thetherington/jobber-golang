Technology Changes:

-   gRPC communication between gateway and microservices

    -   Interceptors for gateawy authentication middleware (JWT)
    -   Interceptors for user information from http session
    -   Shared protobuf gRPC packages, functions, and compilations
    -   Stuctured communication
    -   Server streaming for gateway notifications and messages
    -   Protobuf encoding of RabbitMQ messages

-   Filebeat stdout log collection for docker and kubernetes

Technology Overlays:

-   Chi HTTP Router for Gateway
-   Alex Edwards SCS Session Management with Redis Backend
-   Go Workspaces to make designing with a common package easier
-   SQL compiler for Auth Microservice DB
-   Hexagonal project structure with dependency injection

Technology Limitations:

-   No automatic error bubbling between microservices and gateway
-   No mongoose ORM
-   No out of box library for direct elasticsearch application logging.
-   No official socket.io framework support (needed to use a general websocket implementation)
