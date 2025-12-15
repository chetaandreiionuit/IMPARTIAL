package com.truthweave.presentation.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.BoxScope
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.theme.glassBackground

// [RO] GLASS CARD WRAPPER
// Reutilizabil pentru orice card (Article, Ad, Profile Widget)
@Composable
fun GlassCard(
    modifier: Modifier = Modifier,
    content: @Composable BoxScope.() -> Unit
) {
    Box(
        modifier = modifier
            .glassBackground(shape = RoundedCornerShape(16.dp))
            .padding(16.dp),
        content = content
    )
}
