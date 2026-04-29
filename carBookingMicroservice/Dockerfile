# syntax=docker/dockerfile:1.6
FROM maven:3.9-eclipse-temurin-17 AS builder
WORKDIR /build

# Cache dependencies. `dependency:go-offline` doesn't fully resolve the
# protobuf-maven-plugin (it pulls protoc at codegen time), but it still
# warms the bulk of the dependency tree.
COPY pom.xml ./
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q dependency:go-offline || true

COPY proto ./proto
COPY src ./src
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q -DskipTests package && \
    cp target/*.jar /build/app.jar

# --- runtime ---------------------------------------------------------
# eclipse-temurin:*-jre-alpine has no linux/arm64 manifest. The jammy tag
# is multi-arch (amd64 + arm64) and only ~80 MB larger.
FROM eclipse-temurin:17-jre-jammy AS runtime

RUN groupadd --system app && useradd --system --gid app app
WORKDIR /app

COPY --from=builder /build/app.jar /app/app.jar

USER app:app
EXPOSE 8080

ENV JAVA_OPTS="-XX:MaxRAMPercentage=75.0"
ENTRYPOINT ["sh", "-c", "exec java $JAVA_OPTS -jar /app/app.jar"]
