package com.uni.carbooking.infrastructure.grpc;

import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.money.Money;
import com.uni.currency.grpc.v1.ConvertRequest;
import com.uni.currency.grpc.v1.ConvertResponse;
import com.uni.currency.grpc.v1.CurrencyConverterGrpc;
import io.grpc.ManagedChannel;
import io.grpc.Server;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.stub.StreamObserver;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

/**
 * Spins up an in-process gRPC server with a configurable
 * {@link CurrencyConverterGrpc.CurrencyConverterImplBase} so the adapter is
 * exercised over a real channel + stub, but without the network.
 */
class GrpcCurrencyConverterClientTest {

    private Server server;
    private ManagedChannel channel;
    private CurrencyConverterGrpc.CurrencyConverterImplBase impl;

    @BeforeEach
    void setUp() throws Exception {
        var serverName = InProcessServerBuilder.generateName();
        // Start with a no-op impl; each test swaps it via setImpl().
        impl = new CurrencyConverterGrpc.CurrencyConverterImplBase() {};
        server = InProcessServerBuilder.forName(serverName)
                .directExecutor()
                .addService(new RoutingService(() -> impl))
                .build()
                .start();
        channel = InProcessChannelBuilder.forName(serverName).directExecutor().build();
    }

    @AfterEach
    void tearDown() throws Exception {
        if (channel != null) {
            channel.shutdownNow().awaitTermination(2, TimeUnit.SECONDS);
        }
        if (server != null) {
            server.shutdownNow().awaitTermination(2, TimeUnit.SECONDS);
        }
    }

    private GrpcCurrencyConverterClient newClient(long deadlineMs) {
        return new GrpcCurrencyConverterClient(
                CurrencyConverterGrpc.newBlockingStub(channel),
                deadlineMs);
    }

    @Test
    void converts_money_on_ok_response() {
        impl = new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void convert(ConvertRequest req, StreamObserver<ConvertResponse> resp) {
                assertThat(req.getAmount()).isEqualTo("100.00");
                assertThat(req.getFromCurrency()).isEqualTo("USD");
                assertThat(req.getToCurrency()).isEqualTo("EUR");
                resp.onNext(ConvertResponse.newBuilder()
                        .setAmount("92.39")
                        .setCurrency("EUR")
                        .setSourceAmount("100.00")
                        .setSourceCurrency("USD")
                        .setRateDate("2026-04-28")
                        .build());
                resp.onCompleted();
            }
        };

        var result = newClient(5000)
                .convert(new Money(new BigDecimal("100.00"), "USD"), "EUR");

        assertThat(result.amount()).isEqualByComparingTo("92.39");
        assertThat(result.currency()).isEqualTo("EUR");
    }

    @Test
    void invalid_argument_becomes_currency_conversion_failed() {
        impl = new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void convert(ConvertRequest req, StreamObserver<ConvertResponse> resp) {
                resp.onError(Status.INVALID_ARGUMENT
                        .withDescription("unknown currency: XYZ")
                        .asRuntimeException());
            }
        };

        assertThatThrownBy(() -> newClient(5000)
                .convert(new Money(new BigDecimal("1"), "USD"), "XYZ"))
                .isInstanceOf(CurrencyConversionFailed.class)
                .hasMessageContaining("INVALID_ARGUMENT")
                .hasMessageContaining("unknown currency");
    }

    @Test
    void unavailable_becomes_currency_conversion_failed() {
        impl = new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void convert(ConvertRequest req, StreamObserver<ConvertResponse> resp) {
                resp.onError(Status.UNAVAILABLE
                        .withDescription("server gone")
                        .asRuntimeException());
            }
        };

        assertThatThrownBy(() -> newClient(5000)
                .convert(new Money(new BigDecimal("1"), "USD"), "EUR"))
                .isInstanceOf(CurrencyConversionFailed.class)
                .hasMessageContaining("unavailable");
    }

    @Test
    void deadline_exceeded_becomes_currency_conversion_failed() {
        impl = new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void convert(ConvertRequest req, StreamObserver<ConvertResponse> resp) {
                // Never respond — the deadline must fire.
                try {
                    Thread.sleep(2000);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
        };

        assertThatThrownBy(() -> newClient(50)
                .convert(new Money(new BigDecimal("1"), "USD"), "EUR"))
                .isInstanceOf(CurrencyConversionFailed.class)
                .hasMessageContaining("timed out");
    }

    /**
     * Forwards Convert calls to the test-supplied impl reference, so each test
     * can swap behavior in-place without rebuilding the server.
     */
    private static final class RoutingService extends CurrencyConverterGrpc.CurrencyConverterImplBase {
        private final java.util.function.Supplier<CurrencyConverterGrpc.CurrencyConverterImplBase> delegate;

        RoutingService(java.util.function.Supplier<CurrencyConverterGrpc.CurrencyConverterImplBase> delegate) {
            this.delegate = delegate;
        }

        @Override
        public void convert(ConvertRequest req, StreamObserver<ConvertResponse> resp) {
            delegate.get().convert(req, resp);
        }
    }
}
