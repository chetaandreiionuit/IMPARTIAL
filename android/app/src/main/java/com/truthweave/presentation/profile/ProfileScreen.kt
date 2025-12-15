package com.truthweave.presentation.profile

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowRight
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.drawscope.Fill
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.components.GlassCard
import com.truthweave.presentation.theme.TruthColors
import kotlin.math.cos
import kotlin.math.sin

@Composable
fun ProfileScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(TruthColors.DeepVoidBlack)
            .padding(16.dp)
    ) {
        // Header
        Row(verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(64.dp).background(Color.Gray, CircleShape))
            Spacer(modifier = Modifier.width(16.dp))
            Column {
                Text("Alex Mercer", style = MaterialTheme.typography.headlineSmall, color = TruthColors.TextPrimary)
                Text("Chief Data Architect", style = MaterialTheme.typography.bodyMedium, color = TruthColors.NeonCyan)
            }
        }
        
        Spacer(modifier = Modifier.height(32.dp))
        
        Text("Cognitive Diet", style = MaterialTheme.typography.titleMedium, color = TruthColors.TextSecondary)
        Spacer(modifier = Modifier.height(16.dp))
        
        // Radar Chart
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(300.dp),
            contentAlignment = Alignment.Center
        ) {
            RadarChart()
        }
        
        // Stats
        Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceEvenly) {
            StatItem("60%", "Politics")
            StatItem("40%", "Science")
        }
        
        Spacer(modifier = Modifier.height(32.dp))
        
        // Settings List
        SettingsItem("Manage Sources")
        Spacer(modifier = Modifier.height(8.dp))
        SettingsItem("AI Notification Filter")
        Spacer(modifier = Modifier.height(8.dp))
        SettingsItem("Ghost Protocol")
    }
}

@Composable
fun StatItem(value: String, label: String) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(value, style = MaterialTheme.typography.headlineMedium, color = TruthColors.VerifiedGreen)
        Text(label, style = MaterialTheme.typography.labelSmall, color = TruthColors.TextSecondary)
    }
}

@Composable
fun SettingsItem(text: String) {
    GlassCard(Modifier.fillMaxWidth()) {
        Row(verticalAlignment = Alignment.CenterVertically) {
            Text(text, color = Color.White)
            Spacer(modifier = Modifier.weight(1f))
            Icon(Icons.Default.KeyboardArrowRight, null, tint = Color.Gray)
        }
    }
}

@Composable
fun RadarChart() {
    val labels = listOf("Politics", "Tech", "Econ", "Sci", "Global")
    val values = listOf(0.8f, 0.6f, 0.4f, 0.9f, 0.7f) // 0.0 to 1.0
    
    Canvas(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        val centerX = size.width / 2
        val centerY = size.height / 2
        val radius = minOf(centerX, centerY)
        
        val angleStep = (2 * Math.PI / labels.size).toFloat()
        
        // Draw Grid (Concentric Pentagons)
        for (i in 1..4) {
            val r = radius * (i / 4f)
            val path = Path()
            for (j in labels.indices) {
                val angle = j * angleStep - (Math.PI / 2).toFloat() // Start at top
                val x = centerX + r * cos(angle)
                val y = centerY + r * sin(angle)
                if (j == 0) path.moveTo(x, y) else path.lineTo(x, y)
            }
            path.close()
            drawPath(path, Color.Gray.copy(alpha = 0.3f), style = Stroke(width = 1.dp.toPx()))
        }
        
        // Draw Data Polygon
        val dataPath = Path()
        for (j in labels.indices) {
            val r = radius * values[j]
            val angle = j * angleStep - (Math.PI / 2).toFloat()
            val x = centerX + r * cos(angle)
            val y = centerY + r * sin(angle)
            if (j == 0) dataPath.moveTo(x, y) else dataPath.lineTo(x, y)
            
            // Draw Dots
            drawCircle(TruthColors.NeonCyan, 4.dp.toPx(), Offset(x, y))
        }
        dataPath.close()
        
        drawPath(dataPath, TruthColors.NeonCyan.copy(alpha = 0.2f), style = Fill)
        drawPath(dataPath, TruthColors.NeonCyan, style = Stroke(width = 2.dp.toPx()))
    }
}
