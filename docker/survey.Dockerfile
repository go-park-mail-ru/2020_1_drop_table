FROM dependencies AS builder

WORKDIR /app

CMD /app/survey_service
