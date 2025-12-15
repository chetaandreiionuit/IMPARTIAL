package com.truthweave.presentation.components

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.truthweave.domain.model.CausalEvent

@Composable
fun CausalMetroMap(
    events: List<CausalEvent>,
    modifier: Modifier = Modifier
) {
    val scrollState = rememberLazyListState()
    
    // "Editorial Fintech" Colors
    val gunmetal = Color(0xFF121212)
    val pathColor = Color(0xFF2E2E2E) // Subtle connection lines

    Box(modifier = modifier.background(gunmetal)) {
        // LAYER 1: The Connecting Lines (Drawn BEHIND the list)
        Canvas(modifier = Modifier.matchParentSize()) {
            val strokeWidth = 2.dp.toPx()
            
            // Simple Metro Track Line
            drawLine(
                color = pathColor,
                start = androidx.compose.ui.geometry.Offset(x = 24.dp.toPx(), y = 0f),
                end = androidx.compose.ui.geometry.Offset(x = 24.dp.toPx(), y = size.height),
                strokeWidth = strokeWidth
            )

            events.forEachIndexed { index, event ->
                 event.effects.forEach { childId ->
                     // Logic to find child index and draw Bezier Curve
                     // path.cubicTo(startX, startY, controlX1, controlY1, endX, endY)
                 }
            }
        }

        // LAYER 2: The Content (LazyColumn)
        LazyColumn(
            state = scrollState,
            contentPadding = PaddingValues(vertical = 16.dp),
            modifier = Modifier.padding(start = 48.dp) // Leave space for the "Metro Line"
        ) {
            items(events, key = { it.id }) { event ->
                EventNodeCard(event)
                Spacer(modifier = Modifier.height(24.dp))
            }
        }
    }
}
