package com.inventory.auth.config;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.interfaces.RSAPrivateCrtKey;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import org.springframework.core.env.Profiles;

@Configuration
@EnableConfigurationProperties(RsaKeyProperties.class)
public class RsaKeyConfig {

    private final RsaKeyProperties properties;

    public RsaKeyConfig(RsaKeyProperties properties) {
        this.properties = properties;
    }

    @Bean
    public KeyPair rsaKeyPair(Environment env) {
        String privateKeyPath = properties.privateKeyPath();
        if (privateKeyPath != null && !privateKeyPath.isBlank()) {
            return loadKeyPairFromPem(privateKeyPath);
        }
        if (env.acceptsProfiles(Profiles.of("prod"))) {
            throw new IllegalStateException(
                    "JWT_PRIVATE_KEY_PATH is required when Spring profile 'prod' is active");
        }
        return generateKeyPair();
    }

    @Bean
    public RSAPublicKey rsaPublicKey(KeyPair keyPair) {
        return (RSAPublicKey) keyPair.getPublic();
    }

    @Bean
    public RSAPrivateKey rsaPrivateKey(KeyPair keyPair) {
        return (RSAPrivateKey) keyPair.getPrivate();
    }

    private KeyPair loadKeyPairFromPem(String privateKeyPath) {
        try {
            String pem = Files.readString(Path.of(privateKeyPath));
            String base64 = pem
                    .replaceAll("-----BEGIN [A-Z ]+-----", "")
                    .replaceAll("-----END [A-Z ]+-----", "")
                    .replaceAll("\\s", "");
            byte[] keyBytes = Base64.getDecoder().decode(base64);
            PKCS8EncodedKeySpec spec = new PKCS8EncodedKeySpec(keyBytes);
            KeyFactory kf = KeyFactory.getInstance("RSA");
            RSAPrivateCrtKey rsaCrt = (RSAPrivateCrtKey) kf.generatePrivate(spec);
            RSAPublicKeySpec pubSpec = new RSAPublicKeySpec(
                    rsaCrt.getModulus(), rsaCrt.getPublicExponent());
            RSAPublicKey publicKey = (RSAPublicKey) kf.generatePublic(pubSpec);
            return new KeyPair(publicKey, rsaCrt);
        } catch (IOException e) {
            throw new IllegalStateException(
                    "Failed to read RSA private key from " + privateKeyPath, e);
        } catch (GeneralSecurityException e) {
            throw new IllegalStateException(
                    "Failed to parse RSA key from " + privateKeyPath, e);
        }
    }

    private KeyPair generateKeyPair() {
        try {
            KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA");
            generator.initialize(2048);
            return generator.generateKeyPair();
        } catch (NoSuchAlgorithmException e) {
            throw new IllegalStateException("Failed to generate RSA key pair", e);
        }
    }
}
