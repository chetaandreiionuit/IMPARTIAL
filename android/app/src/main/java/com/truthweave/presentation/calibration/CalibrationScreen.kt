package com.truthweave.presentation.calibration

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.truthweave.presentation.theme.TruthColors

@OptIn(ExperimentalMaterial3Api::class, ExperimentalLayoutApi::class)
@Composable
fun CalibrationScreen(
    onInitialize: () -> Unit
) {
    var truthThreshold by remember { mutableFloatStateOf(0.8f) }
    
    // Sectors State
    val sectors = listOf("Geopolitics", "Global Markets", "Tech", "Energy", "Crypto", "Science", "Defense")
    val selectedSectors = remember { mutableStateListOf("Geopolitics", "Global Markets", "Tech") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Calibrate Your Oracle", color = TruthColors.TextPrimary) },
                navigationIcon = {
                    IconButton(onClick = {}) {
                         Icon(Icons.Default.ArrowBack, "Back", tint = Color.White)
                    }
                },
                actions = {
                    IconButton(onClick = {}) {
                        Icon(Icons.Default.Settings, "Settings", tint = Color.White)
                    }
                },
                colors = TopAppBarDefaults.mediumTopAppBarColors(containerColor = Color.Transparent)
            )
        },
        containerColor = TruthColors.DeepVoidBlack
    ) { padding ->
        Column(
            modifier = Modifier
                .padding(padding)
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(24.dp)
        ) {
            // Section 1: Truth Threshold
            Text("TRUTH THRESHOLD", style = MaterialTheme.typography.labelSmall, color = TruthColors.TextSecondary)
            Spacer(modifier = Modifier.height(16.dp))
            
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    "${(truthThreshold * 100).toInt()}%",
                    style = MaterialTheme.typography.displayMedium,
                    color = TruthColors.TextPrimary
                )
                Spacer(modifier = Modifier.width(16.dp))
                if (truthThreshold > 0.9f) {
                    Badge(containerColor = TruthColors.VerifiedGreen.copy(alpha = 0.2f)) {
                        Text("STRICT MODE", color = TruthColors.VerifiedGreen, modifier = Modifier.padding(4.dp))
                    }
                }
            }

            Slider(
                value = truthThreshold,
                onValueChange = { truthThreshold = it },
                colors = SliderDefaults.colors(
                    thumbColor = Color.White,
                    activeTrackColor = TruthColors.VerifiedGreen,
                    inactiveTrackColor = TruthColors.TextSecondary.copy(alpha = 0.3f)
                )
            )
            
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                Text("ALL SOURCES", style = MaterialTheme.typography.labelSmall, color = Color.Gray)
                Text("TIER 1 ONLY", style = MaterialTheme.typography.labelSmall, color = Color.Gray)
            }

            Spacer(modifier = Modifier.height(48.dp))

            // Section 2: Active Sectors
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text("ACTIVE SECTORS", style = MaterialTheme.typography.labelSmall, color = TruthColors.TextSecondary)
                Spacer(modifier = Modifier.width(8.dp))
                Text("${selectedSectors.size} Selected", style = MaterialTheme.typography.labelSmall, color = TruthColors.NeonCyan)
            }
            Spacer(modifier = Modifier.height(16.dp))

            // FlowRow (Using simple implementation since FlowRow might be experimental)
            FlowRow(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                sectors.forEach { sector ->
                    val isSelected = selectedSectors.contains(sector)
                    FilterChip(
                        selected = isSelected,
                        onClick = {
                            if (isSelected) selectedSectors.remove(sector) else selectedSectors.add(sector)
                        },
                        label = { Text(sector) },
                        colors = FilterChipDefaults.filterChipColors(
                            selectedContainerColor = TruthColors.NeonCyan.copy(alpha = 0.2f),
                            selectedLabelColor = TruthColors.NeonCyan,
                            containerColor = TruthColors.GlassSurface,
                            labelColor = TruthColors.TextSecondary
                        ),
                        border = FilterChipDefaults.filterChipBorder(
                            borderColor = if (isSelected) TruthColors.NeonCyan else Color.Gray.copy(alpha = 0.5f)
                        )
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(64.dp))
            
            Button(
                onClick = onInitialize,
                modifier = Modifier.fillMaxWidth().height(56.dp),
                colors = ButtonDefaults.buttonColors(containerColor = TruthColors.VerifiedGreen),
                shape = MaterialTheme.shapes.medium
            ) {
                Text("Initialize Feed ->", color = Color.Black, style = MaterialTheme.typography.titleMedium)
            }
        }
    }
}
