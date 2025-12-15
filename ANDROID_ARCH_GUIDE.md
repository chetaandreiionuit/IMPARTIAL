# ANDROID ARCHITECTURE GUIDE: THE ORACLE REVEAL

## 1. The Globe (The Gaia Map)
**Tech Stack**: Mapbox Maps SDK (Mobile) or Google Maps Compose (highly stylized).

**Style Guide**: "Tron Legacy" / "Cyberpunk" aesthetic.
- **Base Map**: Dark Mode (almost black oceans, dark grey land).
- **Points**: Pulse animations.
  - **Red**: Fear/Anger (High Intensity)
  - **Green**: Joy/Prosperity
  - **Blue**: Neutral/Science
- **Optimization**: Fetch `/api/v1/oracle/gaia-map`. It returns lightweight JSON. Render points as `SymbolLayers` or `Circles` in Mapbox for 60fps performance even with 10,000 points.

## 2. The Interaction
- **Spin to Filter**: When the user rotates the globe, listen to `onCameraMove`. Get the current viewport bounds (Lat/Lng).
- **Listing**: Show a "floating list" of headlines overlaying the bottom of the map, filtered by the visible region.

## 3. The Mirror of Truth (Article View)
When opening an article:
- **Header**: Large Neural Rewrite (Neutral).
- **Bias Meter**: A visual gauge showing "Left" vs "Right" bias of the *original* source.
- **Counter-Argument**: If the bias was high, display a specific card:
  > **The Mirror of Truth**
  > "While the source claims X, other perspectives suggest Y because Z."

## 4. Oracle Chat (The Floating FAB)
A persistent "Eye" button.
- **Context**: If reading an article, the chat context includes that article.
- **Global**: If on the map, the chat context is "World History".
- **UI**: Streaming text response. Use Typewriter effect.

## 5. Trust Token
For future Blockchain integration:
- Look for `proofs.arweave_url` in the JSON.
- If present, show a Gold "Immutable" Badge.
- On click, open the Arweave Block Explorer to prove the content hasn't changed since ingestion.
