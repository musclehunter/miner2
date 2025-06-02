-- 町テーブルのカラム確認
SET @column_exists_x = (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_name = 'towns'
    AND column_name = 'position_x'
    AND table_schema = 'minerdb'
);

SET @column_exists_y = (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_name = 'towns'
    AND column_name = 'position_y'
    AND table_schema = 'minerdb'
);

-- 座標カラムが存在しない場合は追加
SET @add_columns = CONCAT(
    'ALTER TABLE towns ',
    CASE WHEN @column_exists_x = 0 THEN 'ADD COLUMN position_x INT NOT NULL DEFAULT 0, ' ELSE '' END,
    CASE WHEN @column_exists_y = 0 THEN 'ADD COLUMN position_y INT NOT NULL DEFAULT 0' ELSE '' END
);

-- カラム追加のSQLを実行（必要な場合）
SET @sql = IF(
    (@column_exists_x = 0 OR @column_exists_y = 0),
    @add_columns,
    'SELECT "座標カラムは既に存在します"'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 既存データの座標情報を更新
UPDATE towns SET position_x = 50, position_y = 15 WHERE id = '1'; -- ルナフロスト
UPDATE towns SET position_x = 80, position_y = 30 WHERE id = '2'; -- スカイスパイア
UPDATE towns SET position_x = 85, position_y = 50 WHERE id = '3'; -- シルヴァリオン
UPDATE towns SET position_x = 70, position_y = 70 WHERE id = '4'; -- インゴットリム
UPDATE towns SET position_x = 50, position_y = 80 WHERE id = '5'; -- ザル＝バディア
UPDATE towns SET position_x = 30, position_y = 70 WHERE id = '6'; -- フォグヴェイル
UPDATE towns SET position_x = 15, position_y = 50 WHERE id = '7'; -- キャメロス

-- 町名も更新
UPDATE towns SET name = 'ルナフロスト', description = '北方のルナリス氷原に位置する半地下都市。氷精族が住まい、常に雪に覆われた美しい都市。' WHERE id = '1';
UPDATE towns SET name = 'スカイスパイア', description = '北東のアエリド天空断崖に浮かぶ空中都市。羽人族が住み、浮遠石と鎖で支えられた壁々が特徴。' WHERE id = '2';
UPDATE towns SET name = 'シルヴァリオン', description = '東方のフェイエルフ深緑大森林に位置する螺旋都市。森エルフ族が住み、巨大な神木を中心に発展した。' WHERE id = '3';
UPDATE towns SET name = 'インゴットリム', description = '南東のマグノーム火山地帯に位置する鉄鋼の都市。溶炉ドワーフ族が住み、火山口に建つ鉱山都市。' WHERE id = '4';
UPDATE towns SET name = 'ザル＝バディア', description = '南方の広大な砂漠に位置するオアシス都市。砂漠獣人族（サンドビーストキン）が住み、地下水脈を利用した都市。' WHERE id = '5';
UPDATE towns SET name = 'フォグヴェイル', description = '南西の湿地に位置する霧に覆われた都市。半霊アンデッド（ミレニアル）族が住み、矢立つ石塔群で特徴づけられる。' WHERE id = '6';
UPDATE towns SET name = 'キャメロス', description = '西方の高原城塁地帯に位置する人間の都市。多層城壁と市場が特徴で、交易ハブとして機能している。' WHERE id = '7';
