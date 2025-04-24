export default {
  name: 'HomePage',
  data () {
    return {
      serverCount: 0,
      routeCount: 0,
      upstreamCount: 0,
    }
  },
  template: `
    <div class="flex flex-wrap gap-4">
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ serverCount }}</span> Servers</el-card>
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ routeCount }}</span> Routes</el-card>
      <el-card shadow="always" class="flex-1"><span class="font-bold">{{ upstreamCount }}</span> Upstreams</el-card>
    </div>
  `,
}