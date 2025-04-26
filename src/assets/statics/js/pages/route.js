import { ElMessage, ElMessageBox } from 'element-plus'
import utils from '../utils/index.js'
import upstream from './upstream.js'
import authentication from './authentication.js'

export default {
  name: 'RoutesListPage',
  data () {
    return {
      loading: false,
      server: {},
      upstreams: [],
      authentications: [],
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
    await this.fetchRelationData()
    await this.fetchList()
    this.loading = false
  },
  methods: {
    async fetchRelationData () {
      try {
        const { data } = await axios.get(`/api/v1/servers/${this.$route.params.id}`)
        this.server = data.Data
      } catch (error) {
        ElMessage.error("fetch server data failed")
      }
      try {
        const { data } = await axios.get('/api/v1/upstreams/')
        this.upstreams = data.Data
      } catch (error) {
        ElMessage.error("fetch upstreams data failed")
      }
      try {
        const { data } = await axios.get('/api/v1/authentications/')
        this.authentications = data.Data
      } catch (error) {
        ElMessage.error("fetch authentications data failed")
      }
    },
    async fetchList () {
      try {
        const { data } = await axios.get(`/api/v1/servers/${this.$route.params.id}/routes`)
        this.list = data.Data
      } catch (error) {
        ElMessage.error("fetch routes data failed")
      }
    },
    async handleDelete (row) {
      try {
        await ElMessageBox.confirm(`Are you sure to delete '${row.Name}' route?`, 'Warning', {
          type: 'warning',
        })
      } catch (error) {
        return
      }

      try {
        this.loading = true
        const { data } = await axios.delete(`/api/v1/servers/${this.$route.params.id}/routes/${row.ID}`)
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
    upstreamFormat (row, column, cellValue, index) {
      return `(${cellValue}) ${_.get(_.find(this.upstreams, {ID: cellValue}), "Name", "")}`
    },
    authenticationFormat (row, column, cellValue, index) {
      return `(${cellValue}) ${_.get(_.find(this.authentications, {ID: cellValue}), "Name", "")}`
    },
    handleEdit (row) {
      const data = _.assign({}, row)
      data.Methods = _.filter(_.split(data.Methods, ','), method => !!method)
      if (_.isEmpty(data.Header)) { data.Header = [] }
      this.form.data = data

      this.form.type = 'update'
      this.$refs.form && this.$refs.form.resetFields()
      this.form.showDrawer = true
    },
    handleCreate () {
      this.form.data = {
        Name: '',
        Methods: [],
        Header: [],
        Path: '',
        StripPath: false,
        UpStreamId: null,
        AuthenticationId: null,
        Enable: true,
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
        const payload = _.assign({}, this.form.data)
        payload.Methods = _.join(_.filter(payload.Methods, method => !!method), ',')

        if (this.form.type === 'create') {
          await axios.post(`/api/v1/servers/${this.$route.params.id}/routes`, payload)
        } else {
          await axios.put(`/api/v1/servers/${this.$route.params.id}/routes/${this.form.data.ID}`, payload)
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
  },
  computed: {
    upstreamOptions() {
      return _.map(this.upstreams, (item) => {
        return {
          label: `${item.ID}: ${item.Name}`,
          value: item.ID,
        }
      })
    },
    authenticationOptions() {
      return _.map(this.authentications, (item) => {
        return {
          label: `${item.ID}: ${item.Name}`,
          value: item.ID,
        }
      })
    },
  },
  template: `
    <div v-loading="loading">
      <div class="flex justify-between mb-2 p-2 bg-white">
        <p>Server: <span class="font-bold">{{ server.Name }}</span> {{ server.Host}}:{{ server.Port }}</p>
        <el-button type="primary" @click="handleCreate">Add Route</el-button>
      </div>
      <el-table :data="list" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="Name" />
        <el-table-column prop="Methods" label="Methods" />
        <el-table-column prop="Path" label="Path" />
        <el-table-column prop="StripPath" label="StripPath" />
        <el-table-column prop="Enable" label="Enable" />

        <el-table-column prop="UpStreamId" label="UpStream" width="120" :formatter="upstreamFormat" />
        <el-table-column prop="AuthenticationId" label="Authentication" width="140" :formatter="authenticationFormat" />
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
        <h4 v-if="form.type === 'update'">Update Route #{{ form.data.ID }}</h4>
        <h4 v-else>Create Route</h4>
      </template>
      <template #default>
        <div>
          <el-form ref="form" :model="form.data" :rules="form.rules" label-position="top">
            <el-form-item label="Name" prop="Name" required>
              <el-input v-model="form.data.Name" placeholder="Alias Name" />
            </el-form-item>

            <el-form-item label="Methods" prop="Methods">
              <el-select v-model="form.data.Methods" placeholder="Select methods" multiple clearable>
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="PATCH" value="PATCH" />
                <el-option label="DELETE" value="DELETE" />
              </el-select>
            </el-form-item>

            <el-form-item label="Path" prop="Path" required>
              <el-input v-model="form.data.Path" placeholder="/base_url" />
            </el-form-item>

            <el-form-item label="UpStream" prop="UpStreamId" required>
              <el-select v-model="form.data.UpStreamId" placeholder="Select upstream">
                <el-option v-for="item in upstreamOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>

            <el-form-item label="Basic Auth" prop="AuthenticationId">
              <el-select v-model="form.data.AuthenticationId" placeholder="Select authentication" clearable>
                <el-option v-for="item in authenticationOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>

            <el-form-item label="StripPath" prop="StripPath">
              <el-checkbox v-model="form.data.StripPath" />
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