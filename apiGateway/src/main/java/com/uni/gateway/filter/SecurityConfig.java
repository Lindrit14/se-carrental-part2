package com.uni.gateway.filter;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;
import org.springframework.web.cors.reactive.CorsConfigurationSource;

/**
 * Reactive (WebFlux) security: validates RS256 JWTs using the shared public
 * key. Permits public routes (healthz, login, register, refresh, password
 * reset, GET /api/v1/cars, GET /api/v1/rates, GET /api/v1/convert) and
 * requires authentication everywhere else under /api/v1/**.
 *
 * The CorsConfigurationSource is wired in here (instead of via a separate
 * CorsWebFilter bean) so CORS headers are written on Spring Security's
 * error responses — without this the browser blocks 401 bodies and the
 * frontend can't trigger its refresh-token retry.
 */
@Configuration
@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    SecurityWebFilterChain securityFilterChain(ServerHttpSecurity http,
                                               CorsConfigurationSource corsConfigSource) {
        http
            .csrf(ServerHttpSecurity.CsrfSpec::disable)
            .cors(cors -> cors.configurationSource(corsConfigSource))
            .authorizeExchange(ex -> ex
                .pathMatchers("/healthz", "/readyz", "/actuator/**").permitAll()
                .pathMatchers(HttpMethod.OPTIONS, "/**").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/login").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/register").permitAll()
                .pathMatchers(HttpMethod.POST, "/api/v1/auth/refresh").permitAll()
                .pathMatchers("/api/v1/auth/password/**").permitAll()
                .pathMatchers(HttpMethod.GET, "/api/v1/cars", "/api/v1/cars/**").permitAll()
                .pathMatchers(HttpMethod.GET, "/api/v1/rates", "/api/v1/rates/**").permitAll()
                .pathMatchers(HttpMethod.GET, "/api/v1/convert", "/api/v1/convert/**").permitAll()
                .pathMatchers("/api/v1/**").authenticated()
                .anyExchange().denyAll()
            )
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(jwt -> {}));
        return http.build();
    }
}
