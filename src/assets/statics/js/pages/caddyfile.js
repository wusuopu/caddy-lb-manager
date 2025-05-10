import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'CaddyfilePage',
  data () {
    return {
      loading: false,
      content: "",
      error: "",
    }
  },
  async mounted () {
    await this.fetchConfig()
  },
  methods: {
    async fetchConfig () {
      this.loading = true
      try {
        const { data } = await axios.get('/api/v1/caddy/config')
        this.content = data.Data
        this.error = ""
      } catch (error) {
        this.error = error.response.data.Error
        ElMessage.error("fetch Caddyfile failed")
      }
      this.loading = false
    },
    async handleReload () {
      try {
        await ElMessageBox.confirm(`Are you sure to reload caddy config`, 'Warning', {
          type: 'warning',
        })
      } catch (error) {
        return
      }

      this.loading = true
      try {
        await axios.post('/api/v1/caddy/reload')
        this.error = ""
        ElMessage.success("Caddyfile has reloaded")
      } catch (error) {
        this.error = error.response.data.Error
        ElMessage.error("reload Caddyfile failed")
      }
      this.loading = false
    },
  },
  template: `
    <div v-loading="loading">
      <div class="flex justify-end gap-2 mb-2 p-2 bg-white">
        <el-button type="primary" @click="fetchConfig">Refresh</el-button>
        <el-button type="danger" @click="handleReload">Caddy Reload Config</el-button>
      </div>

      <div class="p-2 bg-white">
        <p class="text-red-500 mb-2">{{ error }}</p>
        <el-input v-model="content" disabled type="textarea" class="w-full" :rows="30" />
      </div>
    </div>
  `,
}