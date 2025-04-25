import { ElMessage, ElMessageBox } from 'element-plus'
import utils from '../utils/index.js'

export default {
  name: 'CertificateListPage',
  data () {
    return {
      loading: false,
      list: [],
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
        const { data } = await axios.get('/api/v1/caddy/certificates')
        this.list = data.Data
      } catch (error) {
        ElMessage.error("fetch certificates data failed")
      }
    },
    datetimeFormat (row, column, cellValue, index) {
      return utils.datetimeFormat(cellValue)
    },
  },
  template: `
    <div v-loading="loading">
      <el-table :data="list" border stripe style="width: 100%">
        <el-table-column prop="Name" label="Name" width="280" />
        <el-table-column prop="Sans" label="Sans" width="200" />
        <el-table-column prop="ValidityStart" label="ValidityStart" :formatter="datetimeFormat" />
        <el-table-column prop="ValidityEnd" label="ValidityEnd" :formatter="datetimeFormat" />
      </el-table>
    </div>
  `,
}