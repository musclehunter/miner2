<template>
  <div class="min-h-screen bg-gray-100">
    <AdminHeader />
    
    <div class="py-10">
      <header>
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center">
          <h1 class="text-3xl font-bold leading-tight text-gray-900">
            町管理
          </h1>
          <button
            @click="showCreateModal = true"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            新規町を作成
          </button>
        </div>
      </header>
      <main>
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
          <div class="px-4 py-8 sm:px-0">
            <!-- 読み込み中表示 -->
            <div v-if="loading" class="flex justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>読み込み中...</span>
            </div>
            
            <!-- エラー表示 -->
            <div v-else-if="error" class="bg-red-50 p-4 rounded-md">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-red-800">
                    データの取得に失敗しました
                  </h3>
                  <div class="mt-2 text-sm text-red-700">
                    <p>{{ error }}</p>
                  </div>
                  <div class="mt-4">
                    <button
                      @click="fetchTowns"
                      class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                    >
                      再試行
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 町リスト -->
            <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              <div v-for="town in towns" :key="town.id" class="bg-white overflow-hidden shadow rounded-lg">
                <div class="p-5">
                  <div class="flex items-start justify-between">
                    <div>
                      <h3 class="text-lg font-medium text-gray-900">{{ town.name }}</h3>
                      <p class="mt-1 text-sm text-gray-500">{{ town.description }}</p>
                    </div>
                    <div class="ml-4 flex-shrink-0 flex">
                      <button
                        @click="openEditModal(town)"
                        class="bg-white rounded-md text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                      >
                        <span class="sr-only">編集</span>
                        <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                          <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                        </svg>
                      </button>
                      <button
                        @click="confirmDelete(town)"
                        class="ml-2 bg-white rounded-md text-gray-400 hover:text-red-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                      >
                        <span class="sr-only">削除</span>
                        <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  <div class="mt-4 border-t border-gray-200 pt-4">
                    <p class="text-xs text-gray-500">ID: {{ town.id }}</p>
                    <p class="text-xs text-gray-500">作成日: {{ formatDate(town.created_at) }}</p>
                    <p class="text-xs text-gray-500">更新日: {{ formatDate(town.updated_at) }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- 町作成モーダル -->
    <div v-if="showCreateModal" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showCreateModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  新規町を作成
                </h3>
                <div class="mt-4 space-y-4">
                  <div>
                    <label for="name" class="block text-sm font-medium text-gray-700">町名</label>
                    <div class="mt-1">
                      <input type="text" id="name" v-model="newTown.name" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md">
                    </div>
                  </div>
                  <div>
                    <label for="description" class="block text-sm font-medium text-gray-700">説明</label>
                    <div class="mt-1">
                      <textarea id="description" v-model="newTown.description" rows="3" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"></textarea>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              type="button" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              @click="createTown"
              :disabled="creating"
            >
              {{ creating ? '作成中...' : '作成' }}
            </button>
            <button 
              type="button" 
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              @click="showCreateModal = false"
              :disabled="creating"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 町編集モーダル -->
    <div v-if="showEditModal" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showEditModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  町を編集
                </h3>
                <div class="mt-4 space-y-4">
                  <div>
                    <label for="edit-name" class="block text-sm font-medium text-gray-700">町名</label>
                    <div class="mt-1">
                      <input type="text" id="edit-name" v-model="editingTown.name" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md">
                    </div>
                  </div>
                  <div>
                    <label for="edit-description" class="block text-sm font-medium text-gray-700">説明</label>
                    <div class="mt-1">
                      <textarea id="edit-description" v-model="editingTown.description" rows="3" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"></textarea>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              type="button" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              @click="updateTown"
              :disabled="updating"
            >
              {{ updating ? '更新中...' : '更新' }}
            </button>
            <button 
              type="button" 
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              @click="showEditModal = false"
              :disabled="updating"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 削除確認モーダル -->
    <div v-if="showDeleteModal" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showDeleteModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-red-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  町を削除
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    本当に「{{ deleteTargetTown ? deleteTargetTown.name : '' }}」を削除しますか？この操作は元に戻せません。
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              type="button" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
              @click="deleteTown"
              :disabled="deleting"
            >
              {{ deleting ? '削除中...' : '削除' }}
            </button>
            <button 
              type="button" 
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              @click="showDeleteModal = false"
              :disabled="deleting"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import adminService from '@/services/admin'

export default {
  name: 'AdminTownsView',
  components: {
    AdminHeader
  },
  setup() {
    const towns = ref([])
    const loading = ref(true)
    const error = ref('')
    const showCreateModal = ref(false)
    const showEditModal = ref(false)
    const showDeleteModal = ref(false)
    const newTown = ref({ name: '', description: '' })
    const editingTown = ref({})
    const deleteTargetTown = ref(null)
    const creating = ref(false)
    const updating = ref(false)
    const deleting = ref(false)

    // 町一覧取得
    const fetchTowns = async () => {
      loading.value = true
      error.value = ''
      
      try {
        towns.value = await adminService.getAllTowns()
      } catch (err) {
        console.error('町一覧取得エラー:', err)
        error.value = '町情報の取得に失敗しました'
      } finally {
        loading.value = false
      }
    }

    // 日付フォーマット
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // 町作成
    const createTown = async () => {
      if (!newTown.value.name || !newTown.value.description) {
        alert('町名と説明を入力してください')
        return
      }
      
      creating.value = true
      
      try {
        await adminService.createTown(newTown.value)
        
        // 町一覧を再取得
        await fetchTowns()
        
        // モーダルを閉じて入力をリセット
        showCreateModal.value = false
        newTown.value = { name: '', description: '' }
      } catch (err) {
        console.error('町作成エラー:', err)
        alert('町の作成に失敗しました')
      } finally {
        creating.value = false
      }
    }

    // 編集モーダル表示
    const openEditModal = (town) => {
      editingTown.value = { ...town }
      showEditModal.value = true
    }

    // 町更新
    const updateTown = async () => {
      if (!editingTown.value.name || !editingTown.value.description) {
        alert('町名と説明を入力してください')
        return
      }
      
      updating.value = true
      
      try {
        await adminService.updateTown(editingTown.value.id, {
          name: editingTown.value.name,
          description: editingTown.value.description
        })
        
        // 町一覧を再取得
        await fetchTowns()
        
        // モーダルを閉じる
        showEditModal.value = false
      } catch (err) {
        console.error('町更新エラー:', err)
        alert('町の更新に失敗しました')
      } finally {
        updating.value = false
      }
    }

    // 削除確認
    const confirmDelete = (town) => {
      deleteTargetTown.value = town
      showDeleteModal.value = true
    }

    // 町削除
    const deleteTown = async () => {
      if (!deleteTargetTown.value) return
      
      deleting.value = true
      
      try {
        await adminService.deleteTown(deleteTargetTown.value.id)
        
        // 町一覧を再取得
        await fetchTowns()
        
        // モーダルを閉じる
        showDeleteModal.value = false
      } catch (err) {
        console.error('町削除エラー:', err)
        alert('町の削除に失敗しました')
      } finally {
        deleting.value = false
      }
    }

    // コンポーネントマウント時にデータを取得
    onMounted(() => {
      fetchTowns()
    })

    return {
      towns,
      loading,
      error,
      showCreateModal,
      showEditModal,
      showDeleteModal,
      newTown,
      editingTown,
      deleteTargetTown,
      creating,
      updating,
      deleting,
      fetchTowns,
      formatDate,
      createTown,
      openEditModal,
      updateTown,
      confirmDelete,
      deleteTown
    }
  }
}
</script>
