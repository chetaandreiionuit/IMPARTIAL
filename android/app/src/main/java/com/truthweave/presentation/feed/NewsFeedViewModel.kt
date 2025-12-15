package com.truthweave.presentation.feed

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.paging.PagingData
import androidx.paging.cachedIn
import com.truthweave.domain.model.FeedItem
import com.truthweave.domain.repository.NewsRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import javax.inject.Inject

// MVI State
data class FeedState(
    // PagingData is handled separately via Flow
    val isRefreshing: Boolean = false
)

// Side-effects (Toasts, Navigation)
sealed class FeedEffect {
    data class ShowError(val message: String) : FeedEffect()
}

// MVI Events (Intents)
sealed class FeedIntent {
    object Refresh : FeedIntent()
    data class ReportError(val message: String) : FeedIntent()
}

@HiltViewModel
class NewsFeedViewModel @Inject constructor(
    private val repo: NewsRepository
) : ViewModel() {

    private val _state = MutableStateFlow(FeedState())
    val state: StateFlow<FeedState> = _state.asStateFlow()

    private val _effect = Channel<FeedEffect>(Channel.BUFFERED)
    val effect = _effect.receiveAsFlow()

    // [RO] Protocol Antigravity: Infinite Scroll Flow cu Caching
    // Fluxul supraviețuiește rotației ecranului.
    val pagingDataFlow: Flow<PagingData<FeedItem>> = repo.getFeed()
        .cachedIn(viewModelScope)

    fun handleIntent(intent: FeedIntent) {
        when (intent) {
            is FeedIntent.Refresh -> refresh()
            is FeedIntent.ReportError -> {
                _effect.trySend(FeedEffect.ShowError(intent.message))
            }
        }
    }

    private fun refresh() {
        // Paging 3 refresh is usually triggered by adapter.refresh() in UI
        // But we can trigger invalidation if needed via repo
    }
}
