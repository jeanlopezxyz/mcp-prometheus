package com.monitoring.prometheus.infrastructure.config;

import io.smallrye.config.ConfigMapping;
import io.smallrye.config.WithDefault;

@ConfigMapping(prefix = "prometheus.api")
public interface PrometheusConfig {

    @WithDefault("http://localhost:9090")
    String url();
}
