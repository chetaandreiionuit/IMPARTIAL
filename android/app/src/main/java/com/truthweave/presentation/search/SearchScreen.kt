package com.truthweave.presentation.search

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.components.GlassCard
import com.truthweave.presentation.theme.TruthColors

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SearchScreen() {
    var query by remember { mutableStateOf("") }
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(TruthColors.DeepVoidBlack)
            .padding(16.dp)
    ) {
        // 1. Terminal Search Bar
        TextField(
            value = query,
            onValueChange = { query = it },
            placeholder = { Text("> Enter query_ parameters...", fontFamily = FontFamily.Monospace) },
            leadingIcon = { Icon(Icons.Default.Search, contentDescription = null, tint = TruthColors.NeonCyan) },
            modifier = Modifier.fillMaxWidth(),
            colors = TextFieldDefaults.colors(
                focusedContainerColor = TruthColors.GlassSurface,
                unfocusedContainerColor = TruthColors.GlassSurface,
                focusedTextColor = TruthColors.NeonCyan,
                unfocusedTextColor = Color.White,
                focusedIndicatorColor = Color.Transparent,
                unfocusedIndicatorColor = Color.Transparent
            ),
            shape = MaterialTheme.shapes.medium
        )

        Spacer(modifier = Modifier.height(16.dp))

        // 2. Filters
        LazyRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            val filters = listOf("Recent", "Date Range", "Entity", "Sentiment", "Causal Depth")
            items(filters.size) { index ->
                SuggestionChip(
                    onClick = {},
                    label = { Text(filters[index]) },
                    colors = SuggestionChipDefaults.suggestionChipColors(
                        containerColor = if(index==0) TruthColors.NeonCyan.copy(alpha=0.2f) else Color.Transparent,
                        labelColor = if(index==0) TruthColors.NeonCyan else TruthColors.TextSecondary
                    ),
                    border = SuggestionChipDefaults.suggestionChipBorder(
                        borderColor = if(index==0) TruthColors.NeonCyan else Color.Gray
                    )
                )
            }
        }

        Spacer(modifier = Modifier.height(24.dp))
        
        Text("Active Clusters", style = MaterialTheme.typography.labelMedium, color = TruthColors.TextTertiary)
        Spacer(modifier = Modifier.height(16.dp))

        // 3. Cluster List
        LazyColumn(verticalArrangement = Arrangement.spacedBy(16.dp)) {
             item { ClusterItem("Global Trade War", "High Confidence", listOf("Semiconductor Shortage", "Market Volatility")) }
             item { ClusterItem("Climate Accords", "Medium Confidence", listOf("carbon_tax_v2", "EV Subsidies")) }
        }
    }
}

@Composable
fun ClusterItem(title: String, confidence: String, subItems: List<String>) {
    Row(modifier = Modifier.fillMaxWidth()) {
        // Visual Bracket
        Box(
            modifier = Modifier
                .width(4.dp)
                .height(120.dp) // Fixed height estimate
                .background(TruthColors.NeonCyan, MaterialTheme.shapes.small)
        )
        Spacer(modifier = Modifier.width(16.dp))
        
        Column {
            GlassCard(Modifier.fillMaxWidth()) {
                Column {
                    Text(title, style = MaterialTheme.typography.titleMedium, color = TruthColors.TextPrimary)
                    Text(confidence, style = MaterialTheme.typography.labelSmall, color = TruthColors.VerifiedGreen)
                }
            }
            
            subItems.forEach { sub ->
                Row(
                    modifier = Modifier
                        .padding(top = 8.dp, start = 16.dp)
                        .fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // L-shape line connector could be drawn here
                    Text("↳ $sub", style = MaterialTheme.typography.bodyMedium, color = TruthColors.TextSecondary)
                    Spacer(modifier = Modifier.weight(1f))
                    Text("->", color = TruthColors.TextTertiary)
                }
            }
        }
    }
}
