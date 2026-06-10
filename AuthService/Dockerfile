FROM maven:3-eclipse-temurin-21 AS base
WORKDIR /app
COPY pom.xml ./
RUN mvn dependency:go-offline -q

FROM base AS development
COPY src ./src
CMD ["mvn", "spring-boot:run"]

FROM base AS builder
COPY src ./src
RUN mvn package -DskipTests -q

FROM eclipse-temurin:21-jre-alpine AS production
RUN apk add --no-cache wget \
    && addgroup -g 10001 -S appgroup \
    && adduser -u 10001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
USER appuser
EXPOSE 8081
ENTRYPOINT ["java", "-jar", "app.jar"]
