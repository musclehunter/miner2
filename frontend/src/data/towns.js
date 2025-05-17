/**
 * 町のモックデータ
 */
export const towns = [
  {
    id: '1',
    name: 'アイアンヒル',
    description: '鉄鉱石の産地として知られる古い鉱山の町。職人たちが多く住み、優れた鍛冶屋がいる。',
    features: [
      '鉄鉱石の取引が盛ん',
      '鍛冶職人が多い',
      '古い鉱山がある'
    ],
    resources: {
      iron: 5,
      copper: 3,
      silver: 1,
      gold: 0,
      crystal: 1
    },
    marketTraits: {
      buyingPower: 4,
      priceVariation: 2,
      specialtyDemand: 'iron'
    },
    population: 6500,
    difficultyLevel: 1,
    position: { x: 30, y: 40 },
    imageUrl: 'iron-hill.svg'
  },
  {
    id: '2',
    name: 'シルバーレイク',
    description: '美しい湖の近くに建つ町。銀鉱石が産出されることで知られ、高級品の取引が盛ん。',
    features: [
      '銀の取引が盛ん',
      '高級品市場がある',
      '美しい湖の景観'
    ],
    resources: {
      iron: 2,
      copper: 2,
      silver: 5,
      gold: 2,
      crystal: 3
    },
    marketTraits: {
      buyingPower: 5,
      priceVariation: 3,
      specialtyDemand: 'silver'
    },
    population: 8200,
    difficultyLevel: 2,
    position: { x: 65, y: 30 },
    imageUrl: 'silver-lake.svg'
  },
  {
    id: '3',
    name: 'ゴールドクレスト',
    description: '山頂に築かれた豪華な町。金鉱脈が発見されて以来、冒険者や富豪が集まる。',
    features: [
      '金の採掘が盛ん',
      '高級市場がある',
      '裕福な顧客が多い'
    ],
    resources: {
      iron: 1,
      copper: 1,
      silver: 3,
      gold: 5,
      crystal: 2
    },
    marketTraits: {
      buyingPower: 5,
      priceVariation: 4,
      specialtyDemand: 'gold'
    },
    population: 4800,
    difficultyLevel: 4,
    position: { x: 80, y: 70 },
    imageUrl: 'gold-crest.svg'
  },
  {
    id: '4',
    name: 'クリスタルヴェイル',
    description: '神秘的な谷に位置する町。特殊な結晶が採れることで知られ、魔術師が多く訪れる。',
    features: [
      '希少な結晶が採取できる',
      '魔術関連の取引が盛ん',
      '神秘的な雰囲気'
    ],
    resources: {
      iron: 1,
      copper: 2,
      silver: 2,
      gold: 3,
      crystal: 5
    },
    marketTraits: {
      buyingPower: 3,
      priceVariation: 5,
      specialtyDemand: 'crystal'
    },
    population: 3200,
    difficultyLevel: 5,
    position: { x: 20, y: 80 },
    imageUrl: 'crystal-vale.svg'
  },
  {
    id: '5',
    name: 'コッパークリーク',
    description: '小さな川沿いに発展した町。銅鉱石が豊富に採れ、工芸品や日用品の製造が盛ん。',
    features: [
      '銅の取引が盛ん',
      '工芸品製造が活発',
      '川沿いの風景'
    ],
    resources: {
      iron: 2,
      copper: 5,
      silver: 1,
      gold: 0,
      crystal: 1
    },
    marketTraits: {
      buyingPower: 3,
      priceVariation: 2,
      specialtyDemand: 'copper'
    },
    population: 5100,
    difficultyLevel: 1,
    position: { x: 50, y: 60 },
    imageUrl: 'copper-creek.svg'
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
