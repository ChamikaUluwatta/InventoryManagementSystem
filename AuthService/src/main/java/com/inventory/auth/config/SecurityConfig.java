package com.inventory.auth.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    @Order(2)
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception{
        http.authorizeHttpRequests((auth) -> auth
        .requestMatchers("/webjars/**", "/css/**").permitAll()
        .requestMatchers("/health").permitAll()
        .requestMatchers("/api/v1/auth/register").permitAll()
        .anyRequest().authenticated())
        .formLogin(Customizer.withDefaults())
        .csrf(csrf -> csrf.disable());
        return http.build();
    }
}
