package com.uni.gateway.filter;

import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.UUID;

/**
 * Ensures every request has an X-Correlation-Id header, generating one if absent.
 * Downstream services log this id so a single trace can be followed across services.
 * Echoed back in the response for client-side correlation.
 */
@Component
public class CorrelationIdFilter implements GlobalFilter, Ordered {

    public static final String HEADER = "X-Correlation-Id";

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        String existing = exchange.getRequest().getHeaders().getFirst(HEADER);
        String correlationId = (existing == null || existing.isBlank()) ? UUID.randomUUID().toString() : existing;

        var mutated = exchange.mutate()
                .request(b -> b.header(HEADER, correlationId))
                .build();
        mutated.getResponse().getHeaders().set(HEADER, correlationId);
        return chain.filter(mutated);
    }

    @Override
    public int getOrder() {
        // Run before everything so all downstream logs include the id.
        return Ordered.HIGHEST_PRECEDENCE;
    }
}
