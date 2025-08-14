FROM openjdk:17-jdk-alpine

WORKDIR /app

# Add a volume pointing to /tmp
VOLUME /tmp

# Copy the executable jar
ARG JAR_FILE=target/*.jar
COPY ${JAR_FILE} app.jar

# Create a non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose port
EXPOSE 8092

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8092/actuator/health || exit 1

# Run the jar file
ENTRYPOINT ["java", "-jar", "/app/app.jar"]