package com.truthweave.presentation.feed

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.paging.compose.LazyPagingItems
import com.truthweave.domain.model.FeedItem
import com.truthweave.presentation.components.CausalTimelineNode
import com.truthweave.presentation.components.GlassCard
import com.truthweave.presentation.theme.TruthColors

@Composable
fun CausalGraphFeed(
    items: LazyPagingItems<FeedItem>,
    modifier: Modifier = Modifier
) {
    // [RO] THE CAUSAL CHAIN FEED
    // Un flux vertical cu o linie continuă și ramificații organice.
    
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(bottom = 100.dp) // Spatiu pentru BottomBar
    ) {
        items(
            count = items.itemCount,
            key = { index -> items[index]?.id ?: index }
        ) { index ->
            val item = items[index]
            
            if (item != null) {
                // Verificăm dacă e "Child Node" (Consequence)
                // În modelul FeedItem actual avem "causalParentId".
                val isChild = (item as? FeedItem.ArticleItem)?.causalParentId != null
                
                CausalFeedRow(
                    item = item,
                    isChild = isChild,
                    isLast = index == items.itemCount - 1
                )
            }
        }
    }
}

@Composable
fun CausalFeedRow(
    item: FeedItem,
    isChild: Boolean,
    isLast: Boolean
) {
    // Row Alignment: [Timeline Graphics] [Content Card]
    IntrinsicHeightRow {
        // 1. The Timeline Graphics Column (Fixed Width)
        CausalTimelineNode(
            isLast = isLast,
            isChild = isChild,
            modifier = Modifier.fillMaxHeight()
        )

        // 2. The Content
        Box(
            modifier = Modifier
                .weight(1f)
                .padding(top = 16.dp, bottom = 0.dp, end = 16.dp)
        ) {
            when (item) {
                is FeedItem.ArticleItem -> ArticleCard(item, isChild)
                // AdItem removed per Editorial Fintech Blueprint
                else -> {} 
            }
        }
    }
}

@Composable
fun IntrinsicHeightRow(content: @Composable RowScope.() -> Unit) {
    Row(
        modifier = Modifier.height(IntrinsicSize.Min),
        content = content
    )
}

@Composable
fun ArticleCard(item: FeedItem.ArticleItem, isChild: Boolean) {
    // Child cards are slightly smaller/indented
    Column {
        if (isChild) {
            Text(
                text = "↳ DIRECT CONSEQUENCE",
                style = MaterialTheme.typography.labelSmall,
                color = TruthColors.NeonCyan,
                modifier = Modifier.padding(bottom = 4.dp)
            )
        }
        
        GlassCard(
             modifier = Modifier.fillMaxWidth()
        ) {
            Column {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    // "99% Verified" Badge
                    Surface(
                         color = if(item.truthScore > 0.8) TruthColors.VerifiedGreen else TruthColors.WarningYellow,
                         shape = androidx.compose.foundation.shape.CircleShape
                    ) {
                        Text(
                            text = "${(item.truthScore * 100).toInt()}% VERIFIED",
                            style = MaterialTheme.typography.labelSmall,
                            color = Color.Black,
                            modifier = Modifier.padding(horizontal = 8.dp, vertical = 2.dp)
                        )
                    }
                    Spacer(modifier = Modifier.weight(1f))
                    Text(
                        text = item.timestamp,
                        style = MaterialTheme.typography.labelSmall,
                        color = TruthColors.TextSecondary
                    )
                }
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Text(
                    text = item.title,
                    style = MaterialTheme.typography.headlineSmall,
                    color = TruthColors.TextPrimary
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Text(
                    text = item.summary,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TruthColors.TextSecondary,
                    maxLines = 3
                )
            }
        }
    }
}
