package com.inventory.auth.controller;

import java.security.interfaces.RSAPublicKey;
import java.util.Map;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.inventory.auth.config.ApiPaths;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.KeyUse;
import com.nimbusds.jose.jwk.RSAKey;

@RestController
public class JwksController {

    private final RSAPublicKey publicKey;

    public JwksController(RSAPublicKey publicKey) {
        this.publicKey = publicKey;
    }

    @GetMapping(ApiPaths.JWKS)
    public Map<String, Object> getJwks() {
        try {
            RSAKey jwk = new RSAKey.Builder(publicKey)
                    .keyUse(KeyUse.SIGNATURE)
                    .algorithm(JWSAlgorithm.RS256)
                    .keyIDFromThumbprint()
                    .build();

            return new JWKSet(jwk).toJSONObject();
        } catch (JOSEException e) {
            throw new RuntimeException("Failed to compute JWK thumbprint", e);
        }
    }
}
