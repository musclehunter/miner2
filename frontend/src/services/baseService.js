import api from './api';

/**
 * 新しい拠点を設立します。
 * @param {string} townId - 拠点を設立する町のID

 * @returns {Promise<any>}
 */
export const createBase = async (townId) => {
  try {
    const response = await api.post('/game/bases', {
      town_id: townId,
    });
    return response.data;
  } catch (error) {
    console.error('Error creating base:', error.response ? error.response.data : error.message);
    throw error;
  }
};
