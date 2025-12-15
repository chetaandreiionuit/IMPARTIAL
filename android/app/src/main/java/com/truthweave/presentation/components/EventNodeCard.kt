package com.truthweave.presentation.components

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.RectangleShape
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.truthweave.domain.model.CausalEvent
import java.time.format.DateTimeFormatter

@Composable
fun EventNodeCard(event: CausalEvent) {
    Card(
        colors = CardDefaults.cardColors(containerColor = Color(0xFF1E2226)), // Glassmorphism surface
        shape = RectangleShape, // Brutalist/Fintech look
        border = BorderStroke(0.5.dp, Color(0xFF333333)),
        modifier = Modifier.fillMaxWidth().padding(end = 16.dp)
    ) {
        Column(modifier = Modifier.padding(12.dp)) {
            // Header: Date & Trust Score
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    text = event.timestamp.format(DateTimeFormatter.ofPattern("HH:mm • dd MMM")),
                    style = MaterialTheme.typography.labelSmall.copy(fontFamily = FontFamily.SansSerif),
                    color = Color.Gray
                )
                Spacer(Modifier.weight(1f))
                // Bridging Consensus Badge
                Text(
                     text = "Trust: ${event.trustScore}",
                     color = Color(0xFF00C853),
                     style = MaterialTheme.typography.labelSmall
                )
            }

            // Body: Neutralized Headline
            Text(
                text = event.neutralSummary, // Use the AI-rewritten title
                style = MaterialTheme.typography.titleMedium.copy(
                    fontFamily = FontFamily.Serif, // SERIF for authority
                    fontWeight = FontWeight.Bold
                ),
                color = Color.White,
                modifier = Modifier.padding(vertical = 8.dp)
            )

            // Footer: Causal Context
            if (event.causes.isNotEmpty()) {
                Text(
                    text = "↳ TRIGGERED BY: ${event.causes.first()}",
                    style = MaterialTheme.typography.bodySmall,
                    color = Color(0xFF2962FF) // Signal Blue
                )
            }
        }
    }
}
