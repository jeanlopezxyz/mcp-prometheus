package com.monitoring.prometheus.infrastructure.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record AlertsResponseDto(
    String status,
    AlertsDataDto data
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record AlertsDataDto(
        List<AlertDto> alerts
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record AlertDto(
        Map<String, String> labels,
        Map<String, String> annotations,
        String state,
        String activeAt,
        String value
    ) {}
}
