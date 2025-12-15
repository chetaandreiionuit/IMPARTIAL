package com.truthweave.data.remote

import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Query

interface TruthWeaveApi {

    @GET("/api/v1/news/feed")
    suspend fun getFeed(
        @Query("page") page: Int,
        @Query("limit") limit: Int
    ): List<FeedDto>

    @POST("/api/v1/chat")
    suspend fun chat(@Body req: ChatRequest): ChatResponse
}

// DTOs
data class FeedDto(
    val type: String, // "article" or "ad"
    val content: Map<String, Any> // We need custom deserializer or generic map for MVP
)

data class ChatRequest(val article_id: String?, val question: String)
data class ChatResponse(val answer: String)
