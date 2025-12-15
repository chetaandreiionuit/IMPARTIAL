package com.truthweave.presentation.theme

import android.graphics.RenderEffect
import android.graphics.Shader
import android.os.Build
import androidx.annotation.RequiresApi
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.*
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

// [RO] ANTIGRAVITY GLASS MODIFIER
// Aplică efectele de sticlă: Blur + Noise (simulat prin Alpha) + Border Luminos
fun Modifier.glassBackground(
    shape: Shape = RoundedCornerShape(16.dp),
    blurRadius: Dp = 20.dp
): Modifier = this
    .clip(shape)
    .background(
        brush = Brush.linearGradient(
            colors = listOf(
                Color.White.copy(alpha = 0.15f), // Sticlă luminoasă sus-stânga
                Color.White.copy(alpha = 0.05f)  // Sticlă mată jos-dreapta
            ),
            start = Offset(0f, 0f),
            end = Offset(Float.POSITIVE_INFINITY, Float.POSITIVE_INFINITY)
        ),
        shape = shape
    )
    .border(
        width = 1.dp,
        brush = Brush.linearGradient(
            colors = listOf(
                Color.White.copy(alpha = 0.40f), // Reflexie puternică sus
                Color.White.copy(alpha = 0.10f)  // Discret jos
            )
        ),
        shape = shape
    )
    .then(
        // Blur-ul real (doar pe Android 12+)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            Modifier.customBlur(blurRadius)
        } else {
            // Fallback: Mai multă opacitate pentru a compensa lipsa blur-ului
            Modifier.background(Color.Black.copy(alpha = 0.3f))
        }
    )

@RequiresApi(Build.VERSION_CODES.S)
private fun Modifier.customBlur(radius: Dp): Modifier {
    return this.graphicsLayer {
        renderEffect = RenderEffect.createBlurEffect(
            radius.toPx(),
            radius.toPx(),
            Shader.TileMode.MIRROR
        ).asComposeRenderEffect()
        
        // Optimizare pentru Antigravity: Clipping la margini
        clip = true 
    }
}
