package com.truthweave.presentation.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color


@Composable
fun TruthWeaveTheme(
    content: @Composable () -> Unit
) {
    val colorScheme = darkColorScheme(
        primary = TruthColors.NeonCyan,
        secondary = TruthColors.VerifiedGreen,
        tertiary = TruthColors.RoyalBlue,
        background = TruthColors.DeepVoidBlack,
        surface = TruthColors.GlassSurface,
        onPrimary = Color.Black,
        onSecondary = Color.Black,
        onTertiary = Color.White,
        onBackground = TruthColors.TextPrimary,
        onSurface = TruthColors.TextPrimary
    )

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}
