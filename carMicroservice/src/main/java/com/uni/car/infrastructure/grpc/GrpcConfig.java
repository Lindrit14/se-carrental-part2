package com.uni.car.infrastructure.grpc;

import com.uni.currency.grpc.v1.CurrencyConverterGrpc;
import io.grpc.ManagedChannel;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import jakarta.annotation.PreDestroy;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.concurrent.TimeUnit;

/**
 * Wires the gRPC channel + blocking stub for the currency-converter service.
 *
 * <p>Address format: {@code host:port}. Plaintext only — TLS would be added at
 * the edge / mesh layer. Per-call deadline is set in the adapter.
 */
@Configuration
public class GrpcConfig {

    private final String address;
    private ManagedChannel channel;

    public GrpcConfig(@Value("${grpc.client.currency-converter.address}") String address) {
        this.address = address;
    }

    @Bean
    public ManagedChannel currencyConverterChannel() {
        var hostPort = parseHostPort(address);
        this.channel = NettyChannelBuilder
                .forAddress(hostPort.host(), hostPort.port())
                .usePlaintext()
                .build();
        return this.channel;
    }

    @Bean
    public CurrencyConverterGrpc.CurrencyConverterBlockingStub currencyConverterStub(
            ManagedChannel channel) {
        return CurrencyConverterGrpc.newBlockingStub(channel);
    }

    @PreDestroy
    public void shutdown() throws InterruptedException {
        if (channel != null) {
            channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
        }
    }

    private record HostPort(String host, int port) {}

    private static HostPort parseHostPort(String raw) {
        var stripped = raw.startsWith("static://") ? raw.substring("static://".length()) : raw;
        var idx = stripped.lastIndexOf(':');
        if (idx <= 0 || idx == stripped.length() - 1) {
            throw new IllegalArgumentException(
                    "grpc.client.currency-converter.address must be host:port, got: " + raw);
        }
        return new HostPort(stripped.substring(0, idx), Integer.parseInt(stripped.substring(idx + 1)));
    }
}
