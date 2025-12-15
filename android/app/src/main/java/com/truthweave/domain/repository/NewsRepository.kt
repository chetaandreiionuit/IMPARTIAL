package com.truthweave.domain.repository

import com.truthweave.domain.model.FeedItem
import kotlinx.coroutines.flow.Flow
import androidx.paging.PagingData

interface NewsRepository {
    // Returns a flow of paging data
    fun getFeed(): Flow<PagingData<FeedItem>>
    
    // Refresh data triggers network fetch
    suspend fun refreshFeed()
    
    // Chat with Oracle
    suspend fun askOracle(articleId: String?, question: String): String
}
