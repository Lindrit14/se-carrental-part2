package com.uni.gateway.filter;

import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.security.core.context.ReactiveSecurityContextHolder;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.List;

/**
 * After Spring Security has validated the JWT, copy useful claims into request
 * headers so downstream services can read them without re-parsing the token.
 *
 * Headers added:
 *   X-User-Id    — JWT 'sub' claim
 *   X-User-Roles — comma-joined 'roles' claim (e.g. "user,admin")
 *
 * Downstream services should still validate signatures themselves until we
 * move to a fully internal-trust model (Phase 5+).
 */
@Component
public class UserClaimsForwardFilter implements GlobalFilter, Ordered {

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        return ReactiveSecurityContextHolder.getContext()
                .map(ctx -> ctx.getAuthentication())
                .filter(auth -> auth instanceof JwtAuthenticationToken)
                .cast(JwtAuthenticationToken.class)
                .map(token -> mutateWithClaims(exchange, token.getToken()))
                .defaultIfEmpty(exchange)
                .flatMap(chain::filter);
    }

    private ServerWebExchange mutateWithClaims(ServerWebExchange exchange, Jwt jwt) {
        String sub = jwt.getSubject();
        List<String> roles = jwt.getClaimAsStringList("roles");
        String rolesHeader = roles == null ? "" : String.join(",", roles);

        return exchange.mutate()
                .request(b -> {
                    if (sub != null) b.header("X-User-Id", sub);
                    b.header("X-User-Roles", rolesHeader);
                })
                .build();
    }

    @Override
    public int getOrder() {
        // Run after security has populated the context but before routing.
        return -1;
    }
}
