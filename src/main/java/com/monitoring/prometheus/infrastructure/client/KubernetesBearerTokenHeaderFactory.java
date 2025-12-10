package com.monitoring.prometheus.infrastructure.client;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.ws.rs.core.MultivaluedHashMap;
import jakarta.ws.rs.core.MultivaluedMap;
import org.eclipse.microprofile.rest.client.ext.ClientHeadersFactory;
import org.jboss.logging.Logger;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

/**
 * ClientHeadersFactory that injects the Kubernetes ServiceAccount bearer token
 * for authentication to OpenShift monitoring stack (Prometheus, AlertManager).
 */
@ApplicationScoped
public class KubernetesBearerTokenHeaderFactory implements ClientHeadersFactory {

    private static final Logger LOG = Logger.getLogger(KubernetesBearerTokenHeaderFactory.class);
    private static final String TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/token";

    @Override
    public MultivaluedMap<String, String> update(MultivaluedMap<String, String> incomingHeaders,
                                                  MultivaluedMap<String, String> clientOutgoingHeaders) {
        MultivaluedMap<String, String> headers = new MultivaluedHashMap<>();

        try {
            Path tokenFile = Path.of(TOKEN_PATH);
            if (Files.exists(tokenFile)) {
                String token = Files.readString(tokenFile).trim();
                headers.add("Authorization", "Bearer " + token);
                LOG.debug("Added ServiceAccount bearer token to request");
            } else {
                LOG.warn("ServiceAccount token not found at " + TOKEN_PATH +
                        " - running outside Kubernetes or token not mounted");
            }
        } catch (IOException e) {
            LOG.error("Failed to read ServiceAccount token: " + e.getMessage());
        }

        return headers;
    }
}
