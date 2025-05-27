import axios from 'axios'

// 管理者認証付きAPIクライアントを作成
const createAuthenticatedClient = () => {
  const adminToken = localStorage.getItem('adminToken')
  
  if (!adminToken) {
    throw new Error('管理者認証が必要です')
  }
  
  // 環境変数からベースURLを取得
  const baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8090'
  
  // ベースURLに「/api」が含まれている場合、重複を避けるために変換
  const adjustedBaseURL = baseURL.endsWith('/api')
    ? baseURL.substring(0, baseURL.length - 4) // "/api"を削除
    : baseURL
  
  console.log('管理者APIクライアントベースURL:', adjustedBaseURL)
  
  return axios.create({
    baseURL: adjustedBaseURL,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    }
  })
}

// 管理者ログイン
const login = async (secretKey) => {
  try {
    // 環境変数からベースURLを取得
    const baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8090'
    
    // ベースURLが「/api」で終わっている場合は、重複を避ける
    const apiUrl = baseURL.endsWith('/api') 
      ? `${baseURL.substring(0, baseURL.length - 4)}/api/admin/login`
      : `${baseURL}/admin/login`
      
    console.log('管理者ログインリクエストURL:', apiUrl)
    
    const response = await axios.post(apiUrl, { secret_key: secretKey })
    return response.data
  } catch (error) {
    console.error('管理者ログインエラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 全ユーザー取得
const getAllUsers = async () => {
  try {
    const client = createAuthenticatedClient()
    console.log('ユーザー一覧取得リクエスト開始')
    const response = await client.get('/api/admin/users')
    console.log('ユーザー一覧取得成功:', response.data)
    return response.data.users || []
  } catch (error) {
    console.error('ユーザー一覧取得エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// ユーザー詳細取得
const getUserDetail = async (userId) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.get(`/api/admin/users/${userId}`)
    return response.data.user
  } catch (error) {
    console.error('ユーザー詳細取得エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// ユーザー更新
const updateUser = async (userId, userData) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.put(`/api/admin/users/${userId}`, userData)
    return response.data
  } catch (error) {
    console.error('ユーザー更新エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// ユーザー削除
const deleteUser = async (userId) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.delete(`/api/admin/users/${userId}`)
    return response.data
  } catch (error) {
    console.error('ユーザー削除エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 未確認ユーザー一覧取得
const getAllPendingUsers = async () => {
  try {
    const client = createAuthenticatedClient()
    console.log('未確認ユーザー一覧取得リクエスト開始')
    const response = await client.get('/api/admin/pending-users')
    console.log('未確認ユーザー一覧取得成功:', response.data)
    return response.data.pendingUsers || []
  } catch (error) {
    console.error('未確認ユーザー一覧取得エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 未確認ユーザー削除
const deletePendingUser = async (token) => {
  try {
    const client = createAuthenticatedClient()
    console.log('未確認ユーザー削除リクエスト開始:', token)
    const response = await client.delete(`/api/admin/pending-users/${token}`)
    console.log('未確認ユーザー削除成功:', response.data)
    return response.data
  } catch (error) {
    console.error('未確認ユーザー削除エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 確認メール再送信
const resendVerification = async (email) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.post('/api/auth/resend-verification', { email })
    return response.data
  } catch (error) {
    console.error('確認メール再送信エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 町一覧取得
const getAllTowns = async () => {
  try {
    const client = createAuthenticatedClient()
    console.log('町一覧取得リクエスト開始')
    const response = await client.get('/api/admin/towns')
    console.log('町一覧取得成功:', response.data)
    return response.data.towns || []
  } catch (error) {
    console.error('町一覧取得エラー:', error, error.response ? error.response.data : error.message)
    throw error
  }
}

// 町作成
const createTown = async (townData) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.post('/api/admin/towns', townData)
    return response.data
  } catch (error) {
    console.error('町作成エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 町更新
const updateTown = async (townId, townData) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.put(`/api/admin/towns/${townId}`, townData)
    return response.data
  } catch (error) {
    console.error('町更新エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

// 町削除
const deleteTown = async (townId) => {
  try {
    const client = createAuthenticatedClient()
    const response = await client.delete(`/api/admin/towns/${townId}`)
    return response.data
  } catch (error) {
    console.error('町削除エラー:', error.response ? error.response.data : error.message)
    throw error
  }
}

export default {
  login,
  getAllUsers,
  getUserDetail,
  updateUser,
  deleteUser,
  getAllPendingUsers,
  deletePendingUser,
  resendVerification,
  getAllTowns,
  createTown,
  updateTown,
  deleteTown
}
