import { ElMessage, ElMessageBox } from 'element-plus'
import utils from '../utils/index.js'

export default {
  name: 'UpstreamListPage',
  data () {
    return {
      loading: false,
      list: [],
      form: {
        showDrawer: false,
        loading: false,
        data: {},
        rules: {},
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
        const { data } = await axios.get('/api/v1/upstreams/')
        this.list = data.Data
      } catch (error) {
        ElMessage.error("fetch upstreams data failed")
      }
    },
    async handleDelete (row) {
      try {
        await ElMessageBox.confirm(`Are you sure to delete '${row.Name}' upstream?`, 'Warning', {
          type: 'warning',
        })
      } catch (error) {
        return
      }

      try {
        this.loading = true
        await axios.delete(`/api/v1/upstreams/${row.ID}`)
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
      this.form.data = data

      this.form.type = 'update'
      this.$refs.form && this.$refs.form.resetFields()
      this.form.showDrawer = true
    },
    handleCreate () {
      this.form.data = {
        Name: '',
        Scheme: 'http://',
        Address: '',
      }
      this.form.type = 'create'
      this.$refs.form && this.$refs.form.resetFields()
      this.form.showDrawer = true
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
          await axios.post('/api/v1/upstreams/', this.form.data)
        } else {
          await axios.put(`/api/v1/upstreams/${this.form.data.ID}`, this.form.data)
        }
      } catch (error) {
        ElMessage.error(error.response.data.Error)
        return
      } finally {
        this.form.loading = false
      }

      this.form.showDrawer = false
      await this.fetchList()
    },
  },
  template: `
    <div v-loading="loading">
      <div class="flex justify-end mb-2 p-2 bg-white">
        <el-button type="primary" @click="handleCreate">Add upstream</el-button>
      </div>
      <el-table :data="list" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="Name" width="100" />
        <el-table-column prop="Scheme" label="Scheme" width="100" />
        <el-table-column prop="Address" label="Address" width="320" />
        <el-table-column prop="CreatedAt" label="CreatedAt" :formatter="datetimeFormat" />
        <el-table-column prop="UpdatedAt" label="UpdatedAt" :formatter="datetimeFormat" />
        <el-table-column label="Operation" fixed="right" width="180">
          <template #default="scope">
            <el-button type="danger" @click="handleDelete(scope.row)">Delete</el-button>
            <el-button type="primary" @click="handleEdit(scope.row)">Edit</el-button>
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
            <el-form-item label="Name" prop="Name" required>
              <el-input v-model="form.data.Name" />
            </el-form-item>

            <el-form-item label="Scheme" prop="Scheme" required>
              <el-select v-model="form.data.Scheme" placeholder="Select scheme">
                <el-option label="http://" value="http://" />
                <el-option label="https://" value="https://" />
                <el-option label="h2c://" value="h2c://" />
                <el-option label="unix/" value="unix/" />
              </el-select>
            </el-form-item>

            <el-form-item label="Address" prop="Address" required>
              <el-input v-model="form.data.Address" />
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