package com.uni.gateway.filter;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

/**
 * Reactive (WebFlux) security: validates RS256 JWTs using the shared public key.
 * Permits public routes (healthz, login, register, refresh, password reset,
 * GET /api/v1/cars) and requires authentication everywhere else under /api/v1/**.
 *
 * The actual JWT decoder is built from spring.security.oauth2.resourceserver.jwt.public-key-location
 * in application.yml.
 */
@Configuration
@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    SecurityWebFilterChain securityFilterChain(ServerHttpSecurity http) {
        http
            .csrf(ServerHttpSecurity.CsrfSpec::disable)
            .cors(cors -> {})
            .authorizeExchange(ex -> ex
                .pathMatchers("/healthz", "/readyz", "/actuator/**").permitAll()
                .pathMatchers(HttpMethod.OPTIONS, "/**").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/login").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/register").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/refresh").permitAll()
                .pathMatchers("/api/v1/auth/password/**").permitAll()
                .pathMatchers(HttpMethod.GET, "/api/v1/cars", "/api/v1/cars/**").permitAll()
                .pathMatchers("/api/v1/**").authenticated()
                .anyExchange().denyAll()
            )
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(jwt -> {}));
        return http.build();
    }
}
