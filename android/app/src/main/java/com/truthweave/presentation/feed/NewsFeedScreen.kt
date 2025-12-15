package com.truthweave.presentation.feed

import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.paging.LoadState
import androidx.paging.compose.collectAsLazyPagingItems
import androidx.paging.compose.itemContentType
import androidx.paging.compose.itemKey
import com.truthweave.domain.model.FeedItem
import com.truthweave.presentation.theme.glassBackground

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NewsFeedScreen(
    viewModel: NewsFeedViewModel = hiltViewModel()
) {
    val pagingItems = viewModel.pagingDataFlow.collectAsLazyPagingItems()

    // [RO] Antigravity Background: Gunmetal (Blueprint Requirement)
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFF121212)) // Gunmetal
    ) {

        Scaffold(
            // Facem TopBar transparent pentru a vedea fundalul
            containerColor = Color.Transparent,
            topBar = {
                CenterAlignedTopAppBar(
                    title = { 
                        Text(
                            "TruthWeave Omni",
                            style = MaterialTheme.typography.headlineMedium.copy(color = Color.White)
                        ) 
                    },
                    colors = TopAppBarDefaults.centerAlignedTopAppBarColors(
                        containerColor = Color.Transparent,
                        titleContentColor = Color.White
                    )
                )
            }
        ) { padding ->
            Box(modifier = Modifier.padding(padding).fillMaxSize()) {
                
                // [RO] CAUSAL CHAIN IMPLEMENTATION
                // We have replaced the standard linear feed with the new CausalGraphFeed
                // that visualizes the connections between events (The "Git Graph").
                CausalGraphFeed(
                    items = pagingItems,
                    modifier = Modifier.fillMaxSize()
                )
                
                if (pagingItems.loadState.refresh is LoadState.Loading) {
                     CircularProgressIndicator(
                         modifier = Modifier.align(Alignment.Center),
                         color = Color.Cyan
                     )
                }
            }
        }
    }
}

// [RO] Wrapper animat pentru fiecare rând
@Composable
fun AntigravityItem(index: Int, content: @Composable () -> Unit) {
    // Animăm doar la prima afișare.
    val visibleState = remember { MutableTransitionState(false).apply { targetState = true } }
    
    // Folosim un delay incremental subtil bazat pe index pentru efectul de cascadă
    // Nota: index % 5 pentru a nu avea delay infinit la scroll adânc.
    val enterDelay = (index % 5) * 50 

    androidx.compose.animation.AnimatedVisibility(
        visibleState = visibleState,
        enter = androidx.compose.animation.fadeIn(
            animationSpec = tween(durationMillis = 500, delayMillis = enterDelay)
        ) + androidx.compose.animation.slideInVertically(
            initialOffsetY = { 100 }, // Vin de jos
            animationSpec = spring(dampingRatio = 0.6f, stiffness = Spring.StiffnessLow)
        ) + androidx.compose.animation.scaleIn(
            initialScale = 0.9f,
            animationSpec = spring(dampingRatio = 0.5f, stiffness = Spring.StiffnessMedium)
        ),
        modifier = Modifier.fillMaxWidth()
    ) {
        content()
    }
}

@Composable
fun ArticleGlassPanel(item: FeedItem.ArticleItem) {
    // Înlocuim Card cu Box + GlassModifier
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .glassBackground() // Efectul nostru custom
            .padding(20.dp) // Padding interior
    ) {
        Column {
            Text(
                item.title, 
                style = MaterialTheme.typography.titleLarge.copy(color = Color.White)
            )
            Spacer(modifier = Modifier.height(12.dp))
            Text(
                item.summary, 
                style = MaterialTheme.typography.bodyMedium.copy(color = Color(0xFFE2E8F0)) // Gri deschis
            )
            Spacer(modifier = Modifier.height(16.dp))
            Row(verticalAlignment = Alignment.CenterVertically) {
                // Truth Score Badge - Glowing
                Surface(
                    shape = MaterialTheme.shapes.small,
                    color = getTruthColor(item.truthScore).copy(alpha = 0.2f),
                    border = androidx.compose.foundation.BorderStroke(1.dp, getTruthColor(item.truthScore))
                ) {
                    Text(
                        text = " TRUTH: ${(item.truthScore * 100).toInt()}% ",
                        style = MaterialTheme.typography.labelSmall,
                        color = Color.White,
                        modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp)
                    )
                }
                
                Spacer(modifier = Modifier.weight(1f))
                
                Text(
                    text = item.biasRating.uppercase(),
                    style = MaterialTheme.typography.labelSmall,
                    color = Color.LightGray
                )
            }
        }
    }
}


// Helpers
fun getTruthColor(score: Double): Color {
    return when {
        score >= 0.8 -> Color(0xFF00E676) // Neon Green
        score >= 0.5 -> Color(0xFFFFEA00) // Neon Yellow
        else -> Color(0xFFFF3D00)         // Neon Red
    }
}
