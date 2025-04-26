import { ElMessage, ElMessageBox } from 'element-plus'
import utils from '../utils/index.js'

export default {
  name: 'ServerListPage',
  data () {
    return {
      loading: false,
      list: [],
      form: {
        showDrawer: false,
        loading: false,
        data: {},
        rules: {
          Name: [
            { required: true, message: 'Name is required' },
          ],
          Port: [
            { required: true, validator: (rule, value, cb) => {
              if (!value ) {
                cb(new Error('Port is required'))
              }
              if (this.form.data.EnableSSL && value === 80) {
                cb(new Error('Port 80 is not allowed when enable SSL'))
              }
              if (!this.form.data.EnableSSL && value === 443) {
                cb(new Error('Port 443 is not allowed when disable SSL'))
              }
              cb()
            }, trigger: 'blur'},
          ],
          Host: [
            { validator: (rule, value, cb) => {
              if (this.form.data.EnableSSL && !value.trim()) {
                cb(new Error('Host is required when enable SSL'))
              }
              cb()
            }, trigger: 'blur'},
          ],
        },
        type: '', // create or update
      }
    }
  },
  async mounted () {
    this.loading = true
    await this.fetchList()
    this.loading = false
  },
  methods: {
    async fetchList () {
      try {
        const { data } = await axios.get('/api/v1/servers/')
        this.list = data.Data
      } catch (error) {
        ElMessage.error("fetch servers data failed")
      }
    },
    async handleDelete (row) {
      try {
        await ElMessageBox.confirm(`Are you sure to delete '${row.Name}'  server?`, 'Warning', {
          type: 'warning',
        })
      } catch (error) {
        return
      }

      try {
        this.loading = true
        await axios.delete(`/api/v1/servers/${row.ID}`)
        ElMessage.success('config has changed. please goto Caddyfile page reload caddy config')
      } catch (error) {
        ElMessage.error(error.response.data.Error)
        return
      } finally {
        this.loading = false
      }

      await this.fetchList()
    },
    datetimeFormat (row, column, cellValue, index) {
      return utils.datetimeFormat(cellValue)
    },
    handleEdit (row) {
      const data = _.assign({}, row)
      data.Port = Number(data.Port)
      this.form.data = data

      this.form.type = 'update'
      this.$refs.form && this.$refs.form.resetFields()
      this.form.showDrawer = true
    },
    handleCreate () {
      this.form.data = {
        Name: '',
        Host: '',
        Port: 443,
        EnableSSL: true,
        Enable: true,
      }
      this.form.type = 'create'
      this.$refs.form && this.$refs.form.resetFields()
      this.form.showDrawer = true
    },
    handleSSLChange (value) {
      if (value) {
        this.form.data.Port = 443
      } else {
        this.form.data.Port = 80
      }
    },
    async handleSubmit () {
      try {
        await this.$refs.form.validate()
      } catch (error) {
        return
      }

      try {
        this.form.loading = true

        if (this.form.type === 'create') {
          await axios.post('/api/v1/servers/', this.form.data)
        } else {
          await axios.put(`/api/v1/servers/${this.form.data.ID}`, this.form.data)
        }
        ElMessage.success('config has changed. please goto Caddyfile page reload caddy config')
      } catch (error) {
        ElMessage.error(error.response.data.Error)
        return
      } finally {
        this.form.loading = false
      }

      this.form.showDrawer = false
      await this.fetchList()
    },
    gotoRoutes (server) {
      this.$router.push(`/servers/${server.ID}/routes`)
    },
  },
  template: `
    <div v-loading="loading">
      <div class="flex justify-end mb-2 p-2 bg-white">
        <el-button type="primary" @click="handleCreate">Add Server</el-button>
      </div>
      <el-table :data="list" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="Name" width="100" />
        <el-table-column prop="Host" label="Host" width="220" />
        <el-table-column prop="Port" label="Port" />
        <el-table-column prop="EnableSSL" label="EnableSSL" width="100" />
        <el-table-column prop="Enable" label="Enable" width="100" />
        <el-table-column prop="CreatedAt" label="CreatedAt" :formatter="datetimeFormat" />
        <el-table-column prop="UpdatedAt" label="UpdatedAt" :formatter="datetimeFormat" />
        <el-table-column label="Operation" fixed="right" width="280">
          <template #default="scope">
            <el-button type="danger" @click="handleDelete(scope.row)">Delete</el-button>
            <el-button type="primary" @click="handleEdit(scope.row)">Edit</el-button>
            <el-button @click="gotoRoutes(scope.row)">Routes</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer v-model="form.showDrawer" direction="rtl">
      <template #header>
        <h4 v-if="form.type === 'update'">Update server #{{ form.data.ID }}</h4>
        <h4 v-else>Create server</h4>
      </template>
      <template #default>
        <div>
          <el-form ref="form" :model="form.data" :rules="form.rules" label-position="top">
            <el-form-item label="Name" prop="Name">
              <el-input v-model="form.data.Name" placeholder="Alias Name" />
            </el-form-item>

            <el-form-item label="Host" prop="Host">
              <el-input v-model="form.data.Host" placeholder="example.com" />
            </el-form-item>

            <el-form-item label="Port" prop="Port">
              <el-input-number v-model="form.data.Port" :min="1" :max="65535" controls-position="right" class="w-full" />
            </el-form-item>

            <el-form-item label="EnableSSL" prop="EnableSSL">
              <el-checkbox v-model="form.data.EnableSSL" @change="handleSSLChange" />
            </el-form-item>

            <el-form-item label="Enable" prop="Enable">
              <el-checkbox v-model="form.data.Enable" />
            </el-form-item>
          </el-form>
        </div>
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="form.showDrawer = false">Cancel</el-button>
          <el-button v-loading.fullscreen.lock="form.loading" type="primary" @click="handleSubmit">
            {{ form.type == 'create' ? 'Create' : 'Update'}}
          </el-button>
        </div>
      </template>
    </el-drawer>
  `,
}