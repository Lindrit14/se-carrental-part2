package com.uni.gateway.filter;

import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import reactor.core.publisher.Mono;

/**
 * Rate-limit key by client IP. Uses X-Forwarded-For when present (frontend
 * may sit behind nginx), otherwise the direct remote address.
 */
@Configuration
public class RateLimitConfig {

    @Bean
    KeyResolver ipKeyResolver() {
        return exchange -> {
            String forwarded = exchange.getRequest().getHeaders().getFirst("X-Forwarded-For");
            if (forwarded != null && !forwarded.isBlank()) {
                return Mono.just(forwarded.split(",")[0].trim());
            }
            var addr = exchange.getRequest().getRemoteAddress();
            return Mono.just(addr == null ? "unknown" : addr.getAddress().getHostAddress());
        };
    }
}
