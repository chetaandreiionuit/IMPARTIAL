package com.truthweave.data.repository

import com.truthweave.data.remote.TruthWeaveApi
import com.truthweave.domain.model.FeedItem
import com.truthweave.domain.repository.NewsRepository
import androidx.paging.Pager
import androidx.paging.PagingConfig
import androidx.paging.PagingData
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject

class NewsRepositoryImpl @Inject constructor(
    private val api: TruthWeaveApi
) : NewsRepository {

    override fun getFeed(): Flow<PagingData<FeedItem>> {
        return Pager(
            config = PagingConfig(
                pageSize = 10,
                enablePlaceholders = false
            ),
            pagingSourceFactory = { FeedPagingSource(api) }
        ).flow
    }

    override suspend fun refreshFeed() {
        // In real app, this updates DB, which triggers Flow.
        // Simplified: ViewModel re-calls getFeed().
    }

    override suspend fun askOracle(articleId: String?, question: String): String {
        return api.chat(com.truthweave.data.remote.ChatRequest(articleId, question)).answer
    }
}
