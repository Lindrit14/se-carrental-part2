package com.uni.carbooking.application.port.out;

import java.time.Instant;

public interface Clock {
    Instant now();
}
