package com.truthweave.presentation.onboarding

import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.theme.TruthColors
import kotlin.random.Random

@Composable
fun OnboardingScreen(
    onNavigateToCalibration: () -> Unit
) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.verticalGradient(
                    colors = listOf(TruthColors.GalacticGradientStart, TruthColors.GalacticGradientEnd)
                )
            )
    ) {
        // 1. Particle Background
        ParticleBackground()

        // 2. Main Content
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Spacer(modifier = Modifier.weight(1f))
            
            // The Chain Visual (Simplified CSS-like representation)
            ChainVisual()
            
            Spacer(modifier = Modifier.height(48.dp))

            // Text
            Text(
                text = "Connect",
                style = MaterialTheme.typography.displayLarge,
                color = TruthColors.TextPrimary
            )
            Text(
                text = "the dots.",
                style = MaterialTheme.typography.displayLarge.copy(fontStyle = FontStyle.Italic),
                color = TruthColors.VerifiedGreen
            )

            Spacer(modifier = Modifier.height(16.dp))
            
            Text(
                text = "Discover the hidden causal chains behind global events.",
                style = MaterialTheme.typography.bodyLarge,
                color = TruthColors.TextSecondary,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center
            )

            Spacer(modifier = Modifier.weight(1f))

            // Page Indicator (Static for now)
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Box(Modifier.size(8.dp).background(TruthColors.VerifiedGreen, CircleShape))
                Box(Modifier.size(8.dp).background(Color.Gray.copy(alpha=0.5f), CircleShape))
                Box(Modifier.size(8.dp).background(Color.Gray.copy(alpha=0.5f), CircleShape))
            }

            Spacer(modifier = Modifier.height(32.dp))

            Button(
                onClick = onNavigateToCalibration,
                colors = ButtonDefaults.buttonColors(containerColor = Color.Transparent),
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = "Begin Analysis",
                    color = TruthColors.NeonCyan,
                    style = MaterialTheme.typography.titleMedium
                )
            }
            
            Spacer(modifier = Modifier.height(16.dp))
        }
    }
}

@Composable
fun ParticleBackground() {
    val infiniteTransition = rememberInfiniteTransition(label = "particles")
    val time by infiniteTransition.animateFloat(
        initialValue = 0f,
        targetValue = 1f,
        animationSpec = infiniteRepeatable(
            animation = tween(10000, easing = LinearEasing),
            repeatMode = RepeatMode.Restart
        ),
        label = "time"
    )

    Canvas(modifier = Modifier.fillMaxSize()) {
        val width = size.width
        val height = size.height
        
        // Draw 50 random particles
        // In a real app, we would preserve state to not re-randomize every frame.
        // For static look + slight shimmer, this is fine or use seeded random.
        val random = Random(12345) 
        
        repeat(30) {
            val cx = random.nextFloat() * width
            val cy = (random.nextFloat() * height + time * 100) % height // Simple vertical drift
            val r = random.nextFloat() * 4.dp.toPx()
            val alpha = random.nextFloat() * 0.3f
            
            drawCircle(
                color = Color.White.copy(alpha = alpha),
                radius = r,
                center = Offset(cx, cy)
            )
        }
    }
}

@Composable
fun ChainVisual() {
    // A stylized vertical chain
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        repeat(5) {
             Box(
                 modifier = Modifier
                     .size(12.dp, 24.dp)
                     .background(Color.Transparent)
                     .border(2.dp, TruthColors.NeonCyan.copy(alpha = 0.5f), CircleShape)
             )
        }
    }
}
