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
RUN apk add --no-cache wget
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
EXPOSE 8081
ENTRYPOINT ["java", "-jar", "app.jar"]
