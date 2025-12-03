package com.monitoring.prometheus.infrastructure.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record QueryResponseDto(
    String status,
    QueryDataDto data,
    String errorType,
    String error
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record QueryDataDto(
        String resultType,
        List<QueryResultDto> result
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record QueryResultDto(
        java.util.Map<String, String> metric,
        List<Object> value,
        List<List<Object>> values
    ) {}
}
