package com.miniexchange.txservice;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
public class TxserviceApplication {
    public static void main(String[] args) {
        SpringApplication.run(TxserviceApplication.class, args);
    }
}
