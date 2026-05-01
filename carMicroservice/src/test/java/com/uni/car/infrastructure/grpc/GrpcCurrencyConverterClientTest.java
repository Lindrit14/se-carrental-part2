package com.uni.car.infrastructure.grpc;

import com.uni.car.domain.error.CurrencyServiceUnavailable;
import com.uni.currency.grpc.v1.CurrencyConverterGrpc;
import com.uni.currency.grpc.v1.GetRatesRequest;
import com.uni.currency.grpc.v1.GetRatesResponse;
import com.uni.currency.grpc.v1.Rate;
import io.grpc.ManagedChannel;
import io.grpc.Server;
import io.grpc.Status;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.stub.StreamObserver;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class GrpcCurrencyConverterClientTest {

    private Server server;
    private ManagedChannel channel;

    @AfterEach
    void tearDown() {
        if (channel != null) channel.shutdownNow();
        if (server != null) server.shutdownNow();
    }

    @Test
    void getRates_returns_parsed_rate_set() throws IOException {
        startServer(new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void getRates(GetRatesRequest request, StreamObserver<GetRatesResponse> obs) {
                obs.onNext(GetRatesResponse.newBuilder()
                        .setBase("EUR")
                        .setRateDate("2026-04-28")
                        .addRates(Rate.newBuilder().setCurrency("USD").setRate("1.0823"))
                        .addRates(Rate.newBuilder().setCurrency("GBP").setRate("0.8634"))
                        .build());
                obs.onCompleted();
            }
        });

        var client = new GrpcCurrencyConverterClient(stub(), 5_000);
        var rs = client.getRates();

        assertThat(rs.base()).isEqualTo("EUR");
        assertThat(rs.rateDate().toString()).isEqualTo("2026-04-28");
        assertThat(rs.rates()).hasSize(2);
        assertThat(rs.rates().get("USD")).isEqualByComparingTo("1.0823");
        assertThat(rs.rates().get("GBP")).isEqualByComparingTo("0.8634");
    }

    @Test
    void unavailable_status_maps_to_CurrencyServiceUnavailable() throws IOException {
        startServer(new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void getRates(GetRatesRequest request, StreamObserver<GetRatesResponse> obs) {
                obs.onError(Status.UNAVAILABLE.withDescription("nope").asRuntimeException());
            }
        });

        var client = new GrpcCurrencyConverterClient(stub(), 5_000);
        assertThatThrownBy(client::getRates)
                .isInstanceOf(CurrencyServiceUnavailable.class)
                .hasMessageContaining("currency-converter unavailable");
    }

    @Test
    void failed_precondition_maps_to_CurrencyServiceUnavailable() throws IOException {
        startServer(new CurrencyConverterGrpc.CurrencyConverterImplBase() {
            @Override
            public void getRates(GetRatesRequest request, StreamObserver<GetRatesResponse> obs) {
                obs.onError(Status.FAILED_PRECONDITION
                        .withDescription("ECB rates not yet loaded").asRuntimeException());
            }
        });

        var client = new GrpcCurrencyConverterClient(stub(), 5_000);
        assertThatThrownBy(client::getRates)
                .isInstanceOf(CurrencyServiceUnavailable.class);
    }

    private void startServer(CurrencyConverterGrpc.CurrencyConverterImplBase service)
            throws IOException {
        var name = "test-" + UUID.randomUUID();
        server = InProcessServerBuilder.forName(name).directExecutor().addService(service).build();
        server.start();
        channel = InProcessChannelBuilder.forName(name).directExecutor().build();
    }

    private CurrencyConverterGrpc.CurrencyConverterBlockingStub stub() {
        return CurrencyConverterGrpc.newBlockingStub(channel);
    }
}
