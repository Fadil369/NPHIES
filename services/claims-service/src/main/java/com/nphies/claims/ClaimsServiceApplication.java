package com.nphies.claims;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.cache.annotation.EnableCaching;

@SpringBootApplication
@EnableKafka
@EnableCaching
public class ClaimsServiceApplication {

    public static void main(String[] args) {
        SpringApplication.run(ClaimsServiceApplication.class, args);
    }
}