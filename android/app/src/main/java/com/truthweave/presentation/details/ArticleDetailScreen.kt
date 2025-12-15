package com.truthweave.presentation.details

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowForward
import androidx.compose.material.icons.filled.Share
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.components.GlassCard
import com.truthweave.presentation.theme.TruthColors

@Composable
fun ArticleDetailScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(TruthColors.DeepVoidBlack)
            .padding(16.dp)
    ) {
        // 1. Causal Context Widget
        CausalContextWidget()
        
        Spacer(modifier = Modifier.height(24.dp))
        
        // 2. Header
        Text(
            "The Global Ripple Effect of Silicon Scarcity",
            style = MaterialTheme.typography.headlineLarge,
            color = TruthColors.TextPrimary
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        // Author
        Row(verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(40.dp).background(Color.Gray, CircleShape))
            Spacer(modifier = Modifier.width(12.dp))
            Column {
                Text("Nexus AI", style = MaterialTheme.typography.labelLarge, color = TruthColors.TextPrimary)
                Text("2 mins ago", style = MaterialTheme.typography.labelSmall, color = TruthColors.TextSecondary)
            }
            Spacer(modifier = Modifier.weight(1f))
            Icon(Icons.Default.Share, null, tint = Color.White)
        }
        
        Spacer(modifier = Modifier.height(24.dp))
        
        // 3. Pull Quote
        Row(Modifier.height(IntrinsicSize.Min)) {
            Box(Modifier.width(4.dp).fillMaxHeight().background(TruthColors.VerifiedGreen))
            Spacer(modifier = Modifier.width(16.dp))
            Text(
                "\"The supply chain is not broken; it is restructuring under new geopolitical physics.\"",
                style = MaterialTheme.typography.titleLarge.copy(fontStyle = FontStyle.Italic),
                color = TruthColors.TextSecondary
            )
        }
        
        Spacer(modifier = Modifier.height(32.dp))
        
        // 4. Line Chart
        Text("Semiconductor Lead Times (Weeks)", style = MaterialTheme.typography.labelSmall, color = TruthColors.TextTertiary)
        Spacer(modifier = Modifier.height(16.dp))
        GlassCard(Modifier.height(200.dp).fillMaxWidth()) {
             SimpleLineChart()
        }
        
        Spacer(modifier = Modifier.weight(1f))
        
        // 5. Sticky Action
        Button(
            onClick = {},
            modifier = Modifier.fillMaxWidth(),
            colors = ButtonDefaults.buttonColors(containerColor = TruthColors.NeonCyan)
        ) {
            Text("See 'Auto Price Hike' Data", color = Color.Black)
            Spacer(modifier = Modifier.width(8.dp))
            Icon(Icons.Default.ArrowForward, null, tint = Color.Black)
        }
    }
}

@Composable
fun CausalContextWidget() {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Left: Cause
        SuggestionChip(onClick={}, label={ Text("Caused by: Supply Chain") }, colors = SuggestionChipDefaults.suggestionChipColors(containerColor = Color.Transparent, labelColor = TruthColors.TextSecondary))
        
        // Center: Current
        SuggestionChip(onClick={}, label={ Text("Viewing: Chip Shortage") }, colors = SuggestionChipDefaults.suggestionChipColors(containerColor = TruthColors.NeonCyan.copy(alpha=0.1f), labelColor = TruthColors.NeonCyan))
        
        // Right: Effect
        Text("Leads to...", color = TruthColors.TextTertiary, style = MaterialTheme.typography.labelSmall)
    }
}

@Composable
fun SimpleLineChart() {
    Canvas(modifier = Modifier.fillMaxSize()) {
        val width = size.width
        val height = size.height
        
        // Data points (normalized 0-1)
        val points = listOf(0.2f, 0.3f, 0.25f, 0.6f, 0.8f, 0.75f, 0.9f)
        
        val stepX = width / (points.size - 1)
        
        val path = Path()
        points.forEachIndexed { index, value ->
            val x = index * stepX
            val y = height - (value * height)
            if (index == 0) path.moveTo(x, y) else path.lineTo(x, y)
            
            // Draw dot
            drawCircle(Color.White, 3.dp.toPx(), Offset(x, y))
        }
        
        drawPath(path, TruthColors.VerifiedGreen, style = Stroke(width = 2.dp.toPx()))
    }
}
