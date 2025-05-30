import { createApp, ref } from 'vue'
import ElementPlus from 'element-plus'
import { createRouter, createWebHashHistory } from 'vue-router'
import HomePage from './pages/home.js'
import ServerListPage from './pages/server.js'
import RoutesListPage from './pages/route.js'
import UpstreamListPage from './pages/upstream.js'
import AuthenticationListPage from './pages/authentication.js'
import CertificateListPage from './pages/certificates.js'
import CaddyfilePage from './pages/caddyfile.js'

const app = createApp({
  data () {
    return window._webui_config
  },
  template: /*html*/`
    <el-container class="min-h-[100vh] mx-auto">
      <el-header height="100px">
        <div class="w-full max-w-[1200px] mx-auto text-[#6e7687] font-bold text-2xl">
          <p>Caddy LoadBalancer Manager</p>

          <el-menu mode="horizontal" :router="true" :default-active="$route.name">
            <el-menu-item index="Home" route="/">Dashboard</el-menu-item>
            <el-menu-item index="Servers" route="/servers">Servers</el-menu-item>
            <el-menu-item index="Upstreams" route="/upstreams">Upstreams</el-menu-item>
            <el-menu-item index="Authentications" route="/authentications">Authentications</el-menu-item>
            <el-menu-item index="Certificates" route="/certificates">SSL Certificates</el-menu-item>
            <el-menu-item index="Caddyfile" route="/caddyfile">Caddyfile</el-menu-item>
          </el-menu>
        </div>
      </el-header>


      <el-main class="bg-[#f5f7fb]">
        <div class="max-w-[1200px] mx-auto text-[#6e7687]">
          <router-view></router-view>
        </div>
      </el-main>

      <el-footer height="30px">
        <div class="flex justify-between max-w-[1200px] mx-auto text-[#6e7687]">
          <span>v{{ version }}</span>
          <span><a href="https://github.com/wusuopu/caddy-lb-manager" target="_blank">Fork me on Github</a></span>
        </div>
      </el-footer>
    </el-container>
  `
})

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: HomePage,
    },
    {
      path: '/servers',
      name: 'Servers',
      component: ServerListPage,
    },
    {
      path: '/servers/:id/routes',
      name: 'Routes',
      component: RoutesListPage,
    },
    {
      path: '/upstreams',
      name: 'Upstreams',
      component: UpstreamListPage,
    },
    {
      path: '/authentications',
      name: 'Authentications',
      component: AuthenticationListPage,
    },
    {
      path: '/certificates',
      name: 'Certificates',
      component: CertificateListPage,
    },
    {
      path: '/caddyfile',
      name: 'Caddyfile',
      component: CaddyfilePage,
    },
    {
      path: '/:pathMatch(.*)*',
      component: {
        template: '<div>not found</div>'
      }
    },
  ],
})

app.use(router)
app.use(ElementPlus)

app.mount('#app')
