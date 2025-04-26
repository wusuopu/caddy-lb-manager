export default {
  name: 'HomePage',
  data () {
    return {
      serverCount: 0,
      routeCount: 0,
      upstreamCount: 0,
    }
  },
  async mounted () {
    const { data } = await axios.get('/api/v1/dashboard')
    this.serverCount = data.Data.serverCount
    this.routeCount = data.Data.routeCount
    this.upstreamCount = data.Data.upstreamCount
  },
  template: `
    <div class="flex flex-wrap gap-4">
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ serverCount }}</span> Servers</el-card>
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ routeCount }}</span> Routes</el-card>
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ upstreamCount }}</span> Upstreams</el-card>
    </div>

    <div class="mt-4">
      <p class="font-bold"> Guide:</p>
      <p>
        1. Add some <router-link to="/upstreams" class="text-blue-500">upstreams</router-link>. <br>
        2. Add some <router-link to="/authentications" class="text-blue-500">authentications</router-link> if you need to enable basic_auth. <br>
        3. Add some <router-link to="/servers" class="text-blue-500">servers</router-link> and routes. <br>
        4. After change config, you need to reload <router-link to="/caddyfile" class="text-blue-500">Caddyfile</router-link> config. <br>
      </p>
    </div>
  `,
}