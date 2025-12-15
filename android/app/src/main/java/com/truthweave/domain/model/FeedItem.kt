package com.truthweave.domain.model

sealed class FeedItem {
    abstract val id: String
    
    data class ArticleItem(
        override val id: String,
        val title: String,
        val summary: String,
        val content: String,
        val imageUrl: String?,
        val truthScore: Double,
        val biasRating: String, // "Left", "Right", "Center"
        val timestamp: String,
        val causalParentId: String? = null,
        val lane: Int = 0
    ) : FeedItem()

    data class AdItem(
        override val id: String,
        val title: String,
        val body: String,
        val mediaUrl: String,
        val targetUrl: String,
        val sponsorLabel: String = "Sponsored"
    ) : FeedItem()
}
