FROM dependencies AS builder

WORKDIR /app

EXPOSE 8084 8084
EXPOSE 8083 8083

CMD /app/staff_service
