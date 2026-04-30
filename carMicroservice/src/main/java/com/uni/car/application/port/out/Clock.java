package com.uni.car.application.port.out;

import java.time.Instant;

public interface Clock {
    Instant now();
}
