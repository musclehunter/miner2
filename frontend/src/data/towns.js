/**
 * 町のモックデータ - world_setting.mdに基づいた7つの地域
 */
export const towns = [
  {
    id: '1',
    name: 'ルナフロスト',
    description: '北方のルナリス氷原に位置する半地下都市。氷精族が住まい、常に雪に覆われた美しい都市。',
    features: [
      '氷の結晶の採掘が盛ん',
      '夜間には特別な力を発揮する',
      '満月時に特別なイベントが発生'
    ],
    resources: {
      iron: 2,
      copper: 1,
      silver: 4,
      gold: 1,
      crystal: 5
    },
    marketTraits: {
      buyingPower: 4,
      priceVariation: 3,
      specialtyDemand: 'crystal'
    },
    population: 5800,
    difficultyLevel: 3,
    position: { x: 50, y: 15 },
    imageUrl: 'luna-frost.svg'
  },
  {
    id: '2',
    name: 'スカイスパイア',
    description: '北東のアエリド天空断崖に浮かぶ空中都市。羽人族が住み、浮遠石と鎖で支えられた壁々が特徴。',
    features: [
      '雲上の古代遺跡がある',
      '空を滑空する技術が発達',
      '雷属性の鉱物が采掘可能'
    ],
    resources: {
      iron: 3,
      copper: 2,
      silver: 3,
      gold: 2,
      crystal: 4
    },
    marketTraits: {
      buyingPower: 5,
      priceVariation: 4,
      specialtyDemand: 'thunder_crystal'
    },
    population: 7200,
    difficultyLevel: 4,
    position: { x: 80, y: 30 },
    imageUrl: 'sky-spire.svg'
  },
  {
    id: '3',
    name: 'シルヴァリオン',
    description: '東方のフェイエルフ深緑大森林に位置する螺旋都市。森エルフ族が住み、巨大な神木を中心に発展した。',
    features: [
      '自然系魔法が発達',
      '希少な植物材が采取可能',
      '树木の高所に展開する立体都市'
    ],
    resources: {
      iron: 1,
      copper: 3,
      silver: 2,
      gold: 2,
      crystal: 4
    },
    marketTraits: {
      buyingPower: 4,
      priceVariation: 3,
      specialtyDemand: 'nature_crystal'
    },
    population: 6800,
    difficultyLevel: 3,
    position: { x: 85, y: 50 },
    imageUrl: 'sylvarion.svg'
  },
  {
    id: '4',
    name: 'インゴットリム',
    description: '南東のマグノーム火山地帯に位置する鉄鋼の都市。溶炉ドワーフ族が住み、火山口に建つ鉱山都市。',
    features: [
      '高品質の武器・防具の製造',
      '耐熱装備が必要な鉱山',
      '古代ドワーフの技術が残る'
    ],
    resources: {
      iron: 5,
      copper: 4,
      silver: 2,
      gold: 3,
      crystal: 1
    },
    marketTraits: {
      buyingPower: 4,
      priceVariation: 2,
      specialtyDemand: 'iron'
    },
    population: 5500,
    difficultyLevel: 4,
    position: { x: 70, y: 70 },
    imageUrl: 'ingot-rim.svg'
  },
  {
    id: '5',
    name: 'ザル＝バディア',
    description: '南方の広大な砂漠に位置するオアシス都市。砂漠獣人族（サンドビーストキン）が住み、地下水脈を利用した都市。',
    features: [
      '砂層に無尽の資源が埋もれている',
      '状況に合わせた狩猟技術が発達',
      '日中は炎熱、夜間は凍結する激しい環境'
    ],
    resources: {
      iron: 3,
      copper: 3,
      silver: 2,
      gold: 4,
      crystal: 2
    },
    marketTraits: {
      buyingPower: 4,
      priceVariation: 5,
      specialtyDemand: 'gold'
    },
    population: 4800,
    difficultyLevel: 4,
    position: { x: 50, y: 80 },
    imageUrl: 'zal-badia.svg'
  },
  {
    id: '6',
    name: 'フォグヴェイル',
    description: '南西の湿地に位置する霧に覆われた都市。半霊アンデッド（ミレニアル）族が住み、矢立つ石塔群で特徴づけられる。',
    features: [
      '霊魂と交わる独自の文化',
      '毒・病気などの状態異常に注意が必要',
      '夏の黒泥底には豪華な保管庫が眠る'
    ],
    resources: {
      iron: 2,
      copper: 3,
      silver: 4,
      gold: 3,
      crystal: 3
    },
    marketTraits: {
      buyingPower: 3,
      priceVariation: 4,
      specialtyDemand: 'silver'
    },
    population: 4200,
    difficultyLevel: 5,
    position: { x: 30, y: 70 },
    imageUrl: 'fogveil.svg'
  },
  {
    id: '7',
    name: 'キャメロス',
    description: '西方の高原城塁地帯に位置する人間の都市。多層城壁と市場が特徴で、交易ハブとして機能している。',
    features: [
      '世界各地からの商人が集まる交易拠点',
      '多様な職業ギルドが存在',
      '初心者向けのクエストが豊富'
    ],
    resources: {
      iron: 3,
      copper: 4,
      silver: 3,
      gold: 2,
      crystal: 2
    },
    marketTraits: {
      buyingPower: 5,
      priceVariation: 3,
      specialtyDemand: 'copper'
    },
    population: 12000,
    difficultyLevel: 1,
    position: { x: 15, y: 50 },
    imageUrl: 'cameloth.svg'
  }
];

/**
 * 町のリソース情報を取得
 */
export const getResourceLabel = (resourceValue) => {
  switch (resourceValue) {
    case 0: return '皆無';
    case 1: return '乏しい';
    case 2: return '少量';
    case 3: return '普通';
    case 4: return '豊富';
    case 5: return '非常に豊富';
    default: return '不明';
  }
};

/**
 * 町の難易度表示
 */
export const getDifficultyLabel = (difficultyValue) => {
  switch (difficultyValue) {
    case 1: return '初心者向け';
    case 2: return '容易';
    case 3: return '普通';
    case 4: return '難しい';
    case 5: return '専門家向け';
    default: return '不明';
  }
};

/**
 * IDから町の情報を取得
 */
export const getTownById = (townId) => {
  return towns.find(town => town.id === townId) || null;
};
