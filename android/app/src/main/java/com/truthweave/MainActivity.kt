package com.truthweave

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.navigation.compose.*
import com.truthweave.presentation.calibration.CalibrationScreen
import com.truthweave.presentation.details.ArticleDetailScreen
import com.truthweave.presentation.feed.NewsFeedScreen
import com.truthweave.presentation.onboarding.OnboardingScreen
import com.truthweave.presentation.profile.ProfileScreen
import com.truthweave.presentation.search.SearchScreen
import com.truthweave.presentation.theme.TruthColors
import com.truthweave.presentation.theme.TruthWeaveTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            TruthWeaveTheme {

                TruthWeaveEntry()
            }
        }
    }
}

@Composable
fun TruthWeaveEntry() {
    val navController = rememberNavController()
    val currentRoute = navController.currentBackStackEntryAsState().value?.destination?.route

    val showBottomBar = currentRoute in listOf("feed", "search", "profile")

    Scaffold(
        bottomBar = {
            if (showBottomBar) {
                NavigationBar(
                    containerColor = TruthColors.GlassSurface,
                    contentColor = TruthColors.NeonCyan
                ) {
                    NavigationBarItem(
                        selected = currentRoute == "feed",
                        onClick = { navController.navigate("feed") },
                        icon = { Icon(Icons.Default.Home, null) },
                        label = { Text("Feed") },
                        colors = NavigationBarItemDefaults.colors(
                            selectedIconColor = TruthColors.NeonCyan,
                            selectedTextColor = TruthColors.NeonCyan,
                            unselectedIconColor = Color.Gray,
                            unselectedTextColor = Color.Gray,
                            indicatorColor = TruthColors.NeonCyan.copy(alpha = 0.2f)
                        )
                    )
                    NavigationBarItem(
                        selected = currentRoute == "search",
                        onClick = { navController.navigate("search") },
                        icon = { Icon(Icons.Default.Search, null) },
                        label = { Text("Oracle") },
                        colors = NavigationBarItemDefaults.colors(
                             selectedIconColor = TruthColors.NeonCyan,
                             selectedTextColor = TruthColors.NeonCyan,
                             unselectedIconColor = Color.Gray,
                             indicatorColor = TruthColors.NeonCyan.copy(alpha = 0.2f)
                        )
                    )
                     NavigationBarItem(
                        selected = currentRoute == "profile",
                        onClick = { navController.navigate("profile") },
                        icon = { Icon(Icons.Default.Person, null) },
                        label = { Text("Profile") },
                        colors = NavigationBarItemDefaults.colors(
                             selectedIconColor = TruthColors.NeonCyan,
                             selectedTextColor = TruthColors.NeonCyan,
                             unselectedIconColor = Color.Gray,
                             indicatorColor = TruthColors.NeonCyan.copy(alpha = 0.2f)
                        )
                    )
                }
            }
        },
        containerColor = TruthColors.DeepVoidBlack
    ) { padding ->
        NavHost(
            navController = navController, 
            startDestination = "onboarding",
            modifier = Modifier.padding(padding)
        ) {
            composable("onboarding") {
                OnboardingScreen(onNavigateToCalibration = { navController.navigate("calibration") })
            }
            composable("calibration") {
                CalibrationScreen(onInitialize = { navController.navigate("feed") })
            }
            composable("feed") {
                NewsFeedScreen() 
            }
            composable("search") {
                SearchScreen()
            }
            composable("profile") {
                ProfileScreen()
            }
            composable("details") {
                ArticleDetailScreen()
            }
        }
    }
}
