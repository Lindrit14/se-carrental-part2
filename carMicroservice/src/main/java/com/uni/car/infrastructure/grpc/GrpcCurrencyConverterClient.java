package com.uni.car.infrastructure.grpc;

import com.uni.car.application.port.out.CurrencyConverterPort;
import com.uni.car.domain.error.CurrencyServiceUnavailable;
import com.uni.currency.grpc.v1.CurrencyConverterGrpc;
import com.uni.currency.grpc.v1.GetRatesRequest;
import com.uni.currency.grpc.v1.GetRatesResponse;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

@Component
class GrpcCurrencyConverterClient implements CurrencyConverterPort {

    private final CurrencyConverterGrpc.CurrencyConverterBlockingStub stub;
    private final long deadlineMs;

    GrpcCurrencyConverterClient(
            CurrencyConverterGrpc.CurrencyConverterBlockingStub stub,
            @Value("${grpc.client.currency-converter.deadline-ms:5000}") long deadlineMs) {
        this.stub = stub;
        this.deadlineMs = deadlineMs;
    }

    @Override
    public RateSet getRates() {
        GetRatesResponse reply;
        try {
            reply = stub
                    .withDeadlineAfter(deadlineMs, TimeUnit.MILLISECONDS)
                    .getRates(GetRatesRequest.getDefaultInstance());
        } catch (StatusRuntimeException e) {
            throw new CurrencyServiceUnavailable(messageFor(e.getStatus()));
        }

        Map<String, BigDecimal> rates = new LinkedHashMap<>(reply.getRatesCount() + 1);
        for (var r : reply.getRatesList()) {
            rates.put(r.getCurrency(), new BigDecimal(r.getRate()));
        }
        return new RateSet(reply.getBase(), LocalDate.parse(reply.getRateDate()), rates);
    }

    private static String messageFor(Status status) {
        var code = status.getCode();
        var description = status.getDescription();
        return switch (code) {
            case DEADLINE_EXCEEDED -> "currency-converter timed out";
            case UNAVAILABLE -> "currency-converter unavailable: "
                    + (description == null ? code.name() : description);
            case FAILED_PRECONDITION -> "currency-converter not ready: "
                    + (description == null ? code.name() : description);
            default -> "currency-converter error: " + code.name()
                    + (description == null ? "" : ": " + description);
        };
    }
}
