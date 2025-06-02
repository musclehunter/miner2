-- 不足している2つの町を追加
INSERT INTO towns (id, name, description, position_x, position_y, created_at, updated_at)
VALUES 
('6', 'フォグヴェイル', '南西の湿地に位置する霧に覆われた都市。半霊アンデッド（ミレニアル）族が住み、矢立つ石塔群で特徴づけられる。', 30, 70, NOW(), NOW()),
('7', 'キャメロス', '西方の高原城塁地帯に位置する人間の都市。多層城壁と市場が特徴で、交易ハブとして機能している。', 15, 50, NOW(), NOW());

-- 全ての町データを確認
SELECT id, name, position_x, position_y FROM towns ORDER BY id;
