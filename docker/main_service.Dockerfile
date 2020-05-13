FROM dependencies AS builder

WORKDIR /app

EXPOSE 8080 8080
EXPOSE 8082 8082

CMD /app/main_service
