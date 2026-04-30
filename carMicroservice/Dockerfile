# syntax=docker/dockerfile:1.6
FROM maven:3.9-eclipse-temurin-17 AS builder
WORKDIR /build

COPY pom.xml ./
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q dependency:go-offline || true

COPY src ./src
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q -DskipTests package && \
    cp target/*.jar /build/app.jar

FROM eclipse-temurin:17-jre-jammy AS runtime

RUN groupadd --system app && useradd --system --gid app app
WORKDIR /app

COPY --from=builder /build/app.jar /app/app.jar

USER app:app
EXPOSE 8082

ENV JAVA_OPTS="-XX:MaxRAMPercentage=75.0"
ENTRYPOINT ["sh", "-c", "exec java $JAVA_OPTS -jar /app/app.jar"]
