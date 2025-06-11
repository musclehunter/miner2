import api from './api';

/**
 * プレイヤーの在庫情報（所持金、鉱石、アイテム）を取得します。
 * @returns {Promise<Object>} 在庫情報を含むPromiseオブジェクト
 */
export const getMyInventory = async () => {
  try {
    const response = await api.get('/game/my/inventory');
    return response.data;
  } catch (error) {
    console.error('Error fetching inventory data:', error);
    throw error;
  }
};
