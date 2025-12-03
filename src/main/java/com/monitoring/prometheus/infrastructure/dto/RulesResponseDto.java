package com.monitoring.prometheus.infrastructure.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record RulesResponseDto(
    String status,
    RulesDataDto data
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record RulesDataDto(
        List<RuleGroupDto> groups
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record RuleGroupDto(
        String name,
        String file,
        List<RuleDto> rules,
        Double interval,
        Double limit,
        String evaluationTime,
        String lastEvaluation
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record RuleDto(
        String state,
        String name,
        String query,
        String duration,
        String keepFiringFor,
        Map<String, String> labels,
        Map<String, String> annotations,
        List<AlertDto> alerts,
        String health,
        String evaluationTime,
        String lastEvaluation,
        String type
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
