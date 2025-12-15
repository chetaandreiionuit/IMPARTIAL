package com.truthweave.data.repository

import androidx.paging.PagingSource
import androidx.paging.PagingState
import com.truthweave.data.remote.TruthWeaveApi
import com.truthweave.domain.model.FeedItem

class FeedPagingSource(
    private val api: TruthWeaveApi
) : PagingSource<Int, FeedItem>() {

    override fun getRefreshKey(state: PagingState<Int, FeedItem>): Int? {
        // [RO] Protocol Antigravity: Calculăm ancora pentru a menține poziția
        return state.anchorPosition?.let { anchorPosition ->
            state.closestPageToPosition(anchorPosition)?.prevKey?.plus(1)
                ?: state.closestPageToPosition(anchorPosition)?.nextKey?.minus(1)
        }
    }

    override suspend fun load(params: LoadParams<Int>): LoadResult<Int, FeedItem> {
        // [RO] DEMO MODE: GENERATING MOCK CAUSAL CHAIN DATA
        // Pentru a verifica vizualizarea "Git Graph", returnăm date hardcodate.
        // În producție, acest bloc ar fi înlocuit de apelul API real.
        
        try {
            val page = params.key ?: 1
            if (page > 1) {
                return LoadResult.Page(data = emptyList(), prevKey = 1, nextKey = null)
            }

            val mockItems = listOf(
                // 1. ROOT EVENT
                FeedItem.ArticleItem(
                    id = "root_1",
                    title = "Federal Reserve Announces 0.50% Rate Hike",
                    summary = "In a bold move to combat inflation, the Fed has raised interest rates, signaling a hawkist stance that reverberates through global markets.",
                    content = "Full content...",
                    imageUrl = null,
                    truthScore = 0.98,
                    biasRating = "Center",
                    timestamp = "2 hrs ago",
                    causalParentId = null,
                    lane = 0
                ),

                // 2. DIRECT CONSEQUENCE (Branching)
                FeedItem.ArticleItem(
                    id = "child_1",
                    title = "Tech Sector Sell-Off: NASDAQ Drops 3.5%",
                    summary = "High-growth tech stocks are hammered as borrowing costs rise. The 'Direct Consequence' of the Fed's decision is immediate market volatility.",
                    content = "Full content...",
                    imageUrl = null,
                    truthScore = 0.92,
                    biasRating = "Market Data",
                    timestamp = "1 hr ago",
                    causalParentId = "root_1", // Links to Root
                    lane = 1 // Visual Lane 1
                ),
                
                // 3. SECONDARY CONSEQUENCE
                FeedItem.ArticleItem(
                    id = "child_2",
                    title = "Crypto Winter Deepens: Bitcoin Below $30k",
                    summary = "Risk assets follow tech stocks downward. Institutional liquidity dries up as yields become attractive elsewhere.",
                    content = "Full content...",
                    imageUrl = null,
                    truthScore = 0.85,
                    biasRating = "Speculative",
                    timestamp = "45 mins ago",
                    causalParentId = "root_1",
                    lane = 1
                ),

                // 4. ANOTHER ROOT (Context)
                FeedItem.ArticleItem(
                    id = "root_2",
                    title = "Global Supply Chain 'Healing' Faster Than Expected",
                    summary = "Shipping container rates have normalized to pre-pandemic levels, suggesting inflationary pressures from logistics are easing.",
                    content = "Full content...",
                    imageUrl = null,
                    truthScore = 0.95,
                    biasRating = "Center",
                    timestamp = "3 hrs ago",
                    causalParentId = null,
                    lane = 0
                ),
                
                 // 5. COMPLEX CHAIN (Grandchild)
                FeedItem.ArticleItem(
                    id = "child_3",
                    title = "Retailers Slash Prices on Electronics",
                    summary = "Overstocked inventories and cheaper shipping lead to massive discounts. A direct downstream effect of the supply chain normalization.",
                    content = "...",
                    imageUrl = null,
                    truthScore = 0.89,
                    biasRating = "Consumer",
                    timestamp = "30 mins ago",
                    causalParentId = "root_2",
                    lane = 1
                )
            )

            return LoadResult.Page(
                data = mockItems,
                prevKey = null,
                nextKey = null // End of list for demo
            )
        } catch (e: Exception) {
            return LoadResult.Error(e)
        }
    }
}
