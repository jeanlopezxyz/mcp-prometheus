package com.monitoring.prometheus.infrastructure.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TargetsResponseDto(
    String status,
    TargetsDataDto data
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record TargetsDataDto(
        List<ActiveTargetDto> activeTargets,
        List<DroppedTargetDto> droppedTargets
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record ActiveTargetDto(
        Map<String, String> discoveredLabels,
        Map<String, String> labels,
        String scrapePool,
        String scrapeUrl,
        String globalUrl,
        String lastError,
        String lastScrape,
        String lastScrapeDuration,
        String health,
        String scrapeInterval,
        String scrapeTimeout
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record DroppedTargetDto(
        Map<String, String> discoveredLabels
    ) {}
}
