-- Seed Ads
INSERT INTO ads (id, type, title, body, media_url, target_url, is_active, priority, created_at)
VALUES 
('d290f1ee-6c54-4b01-90e6-d701748f0851', 'native', 'The Future of AI is Here', 'Discover how AI is transforming finance.', 'https://example.com/ai.jpg', 'https://example.com/ai-report', true, 10, NOW()),
('a1b2c3d4-e5f6-7890-1234-567890abcdef', 'banner', 'Crypto Exchange X', '0% Fees on your first trade.', 'https://example.com/crypto.jpg', 'https://example.com/crypto', true, 8, NOW()),
('11223344-5566-7788-9900-aabbccddeeff', 'native', 'Global Summit 2026', 'Tickets available now for the visionary event.', 'https://example.com/summit.jpg', 'https://example.com/tickets', true, 5, NOW());
