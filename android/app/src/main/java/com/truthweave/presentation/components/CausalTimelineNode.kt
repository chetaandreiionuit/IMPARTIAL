package com.truthweave.presentation.components

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.*
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.PathEffect
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.theme.TruthColors

// [RO] CAUSAL TIMELINE NODE
// Desenează linia punctată verticală și logica de ramificare (Bezier curves).
// Această componentă randează DOAR grafica din stânga, conținutul cardului e separat.

@Composable
fun CausalTimelineNode(
    isLast: Boolean,
    isChild: Boolean,
    modifier: Modifier = Modifier
) {
    val timelineColor = TruthColors.TextSecondary.copy(alpha = 0.3f)
    val branchColor = TruthColors.NeonCyan
    
    Canvas(modifier = modifier.width(48.dp).fillMaxHeight()) {
        val centerX = size.width / 2
        
        // 1. Vertical Dashed Line (The Trunk)
        // Dacă e copil, ramura vine de sus-stânga, deci linia principală continuă oricum.
        if (!isLast) {
            drawLine(
                color = timelineColor,
                start = Offset(centerX, 0f),
                end = Offset(centerX, size.height),
                strokeWidth = 2.dp.toPx(),
                pathEffect = PathEffect.dashPathEffect(floatArrayOf(10f, 10f), 0f)
            )
        } else {
             // Dacă e ultimul, desenăm doar jumătate (până la nod)
             drawLine(
                color = timelineColor,
                start = Offset(centerX, 0f),
                end = Offset(centerX, 24.dp.toPx()), // Până la centrul nodului
                strokeWidth = 2.dp.toPx(),
                pathEffect = PathEffect.dashPathEffect(floatArrayOf(10f, 10f), 0f)
             )
        }

        // 2. The Node (Dot)
        drawCircle(
            color = if (isChild) TruthColors.NeonCyan else TruthColors.TextPrimary,
            radius = 4.dp.toPx(),
            center = Offset(centerX, 24.dp.toPx()) // Aliniat cu top-padding-ul cardului
        )

        // 3. Branching Logic (Bezier Curve)
        if (isChild) {
            val path = Path()
            // Start from the main line (visually above)
            // Simulating a connection from a previous "Parent" node
            // Since we can't draw outside bounds easily, we draw the "Entrance" curve 
            // starting from top-left into center.
            
            val startX = 0f
            val startY = -20f // Off-screen top
            val endX = centerX
            val endY = 24.dp.toPx()
            
            path.moveTo(centerX, 0f) // Connects from vertical line
            path.cubicTo(
                centerX, 12.dp.toPx(),
                centerX + 10f, 12.dp.toPx(),
                centerX + 12.dp.toPx(), 24.dp.toPx()
            )
            
            // For the sake of the "Branching" visual, we actually draw a curve 
            // connecting the main trunk to the content on the RIGHT.
            
            val branchPath = Path()
            branchPath.moveTo(centerX, 24.dp.toPx())
            branchPath.cubicTo(
                 centerX + 10.dp.toPx(), 24.dp.toPx(),
                 size.width - 10.dp.toPx(), 24.dp.toPx(),
                 size.width, 24.dp.toPx()
            )

            drawPath(
                path = branchPath,
                color = branchColor,
                style = Stroke(width = 2.dp.toPx())
            )
        }
    }
}
