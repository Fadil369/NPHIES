package com.nphies.claims;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;

/**
 * NPHIES Claims Service Application
 * 
 * Microservice for claims processing, validation, and audit trails.
 * Integrates with eligibility and benefits services.
 */
@SpringBootApplication
@EnableKafka
@EnableJpaAuditing
public class ClaimsServiceApplication {

    public static void main(String[] args) {
        SpringApplication.run(ClaimsServiceApplication.class, args);
    }
}