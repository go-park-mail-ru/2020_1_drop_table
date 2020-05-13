FROM dependencies AS builder

WORKDIR /app

EXPOSE 8085 8085
EXPOSE 8086 8086

CMD /app/survey_service
