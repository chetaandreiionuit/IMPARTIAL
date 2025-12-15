package com.truthweave.domain.model

import java.time.LocalDateTime

data class CausalEvent(
    val id: String,
    val timestamp: LocalDateTime,
    val neutralSummary: String,
    val trustScore: Double,
    val causes: List<String> = emptyList(),
    val effects: List<String> = emptyList()
)
