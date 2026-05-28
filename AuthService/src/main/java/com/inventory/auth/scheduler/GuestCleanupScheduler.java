package com.inventory.auth.scheduler;

import java.time.Instant;
import java.time.temporal.ChronoUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import com.inventory.auth.repository.UserRepository;

@Component
public class GuestCleanupScheduler {

    private static final Logger log = LoggerFactory.getLogger(GuestCleanupScheduler.class);
    private static final int GUEST_TTL_HOURS = 6;

    private final UserRepository userRepository;

    public GuestCleanupScheduler(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    @Scheduled(cron = "0 0 0 * * *")
    public void cleanupExpiredGuests() {
        Instant cutoff = Instant.now().minus(GUEST_TTL_HOURS, ChronoUnit.HOURS);
        var expiredGuests = userRepository.findByIsGuestTrueAndCreatedAtBefore(cutoff);
        if (expiredGuests.isEmpty()) {
            return;
        }
        log.info("Cleaning up {} expired guest account(s)", expiredGuests.size());
        userRepository.deleteAll(expiredGuests);
    }
}
