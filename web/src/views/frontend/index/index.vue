<!-- +----------------------------------------------------------------------
  | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
  +----------------------------------------------------------------------
  | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
  +----------------------------------------------------------------------
  | Licensed ( https://opensource.org/licenses/MIT )
  +----------------------------------------------------------------------
  | Author: 喜羊羊 <751300685@qq.com>
  +---------------------------------------------------------------------- -->
<template>
  <div class="home-page">
    <!-- Hero Section -->
    <header class="relative pt-12 pb-12 px-6 lg:pt-16 lg:pb-20 overflow-hidden">
      <div class="max-w-7xl mx-auto grid lg:grid-cols-2 gap-12 items-center">
        <!-- 文案 -->
        <div class="relative z-10 text-center lg:text-left">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/50 backdrop-blur-sm border border-white/50 shadow-sm mb-6 animate-float">
            <span class="w-2 h-2 rounded-full bg-clay-success"></span>
            <span class="text-sm font-bold text-clay-muted uppercase tracking-wider">v1.0 正式发布</span>
          </div>
          <h1 class="font-heading font-black text-5xl sm:text-6xl lg:text-7xl leading-[1.1] mb-6 text-clay-foreground">
            开箱即用的<br><span class="clay-text-gradient">开源中后台</span>框架
          </h1>
          <p class="text-lg sm:text-xl text-clay-muted mb-8 leading-relaxed max-w-2xl mx-auto lg:mx-0">
            基于 Vue3 + GoFrame 构建的通用开源中后台管理框架。内置权限管理、代码生成、系统监控等核心模块，开箱即用，快速启动你的业务开发。
          </p>
          <div class="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-6">
            <RouterLink
              to="/diagnosis"
              class="w-full sm:w-auto px-10 py-4 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 text-white text-lg font-bold shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 active:shadow-clay-pressed transition-all duration-300 flex items-center justify-center gap-2"
            >
              <span>免费诊断</span>
              <ArtSvgIcon icon="ri:shield-check-line" class="text-lg" />
            </RouterLink>
            <RouterLink
              to="/pricing"
              class="w-full sm:w-auto px-10 py-4 rounded-full bg-white text-clay-foreground text-lg font-bold shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 transition-all duration-300 flex items-center justify-center gap-2"
            >
              <span>服务价格</span>
              <ArtSvgIcon icon="ri:price-tag-3-line" class="text-lg" />
            </RouterLink>
          </div>
        </div>

        <!-- 右侧模拟界面 -->
        <div class="relative z-10 mt-12 lg:mt-0 perspective-1000">
          <div class="relative w-full aspect-[4/3] rounded-[48px] bg-white/40 backdrop-blur-xl shadow-clay-deep border border-[#d1d9e6]/40 p-6 animate-float-delayed transform rotate-y-12 hover:rotate-y-0 transition-transform duration-700 ease-out">
            <!-- 模拟 Header -->
            <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-2xl bg-white shadow-clay-pressed flex items-center justify-center">
                  <div class="w-6 h-6 rounded-full bg-clay-accent"></div>
                </div>
                <div class="h-4 w-32 bg-gray-200 rounded-full shadow-inner"></div>
              </div>
              <div class="flex gap-3">
                <div class="w-10 h-10 rounded-full bg-white shadow-clay-btn"></div>
                <div class="w-10 h-10 rounded-full bg-white shadow-clay-btn"></div>
              </div>
            </div>
            <!-- 模拟 Content -->
            <div class="grid grid-cols-3 gap-6 h-[calc(100%-80px)]">
              <div class="hidden sm:block col-span-1 bg-white/50 rounded-[32px] shadow-clay-pressed p-4">
                <div class="space-y-4 mt-4">
                  <div v-for="i in 5" :key="i" class="h-8 w-full rounded-xl" :class="i === 2 ? 'bg-clay-accent shadow-clay-btn' : 'bg-transparent'"></div>
                </div>
              </div>
              <div class="col-span-3 sm:col-span-2 space-y-6">
                <div class="grid grid-cols-2 gap-4">
                  <div class="h-32 rounded-[24px] bg-white shadow-clay-card p-4 flex flex-col justify-between">
                    <div class="w-8 h-8 rounded-full bg-blue-100 text-blue-500 flex items-center justify-center">
                      <ArtSvgIcon icon="ri:pie-chart-line" class="text-base" />
                    </div>
                    <div class="h-2 w-16 bg-gray-200 rounded-full"></div>
                    <div class="h-6 w-12 bg-gray-300 rounded-full"></div>
                  </div>
                  <div class="h-32 rounded-[24px] bg-white shadow-clay-card p-4 flex flex-col justify-between">
                    <div class="w-8 h-8 rounded-full bg-green-100 text-green-500 flex items-center justify-center">
                      <ArtSvgIcon icon="ri:bar-chart-box-line" class="text-base" />
                    </div>
                    <div class="h-2 w-16 bg-gray-200 rounded-full"></div>
                    <div class="h-6 w-12 bg-gray-300 rounded-full"></div>
                  </div>
                </div>
                <div class="h-40 rounded-[24px] bg-white shadow-clay-card p-4">
                  <div class="h-4 w-32 bg-gray-200 rounded-full mb-4"></div>
                  <div class="space-y-3">
                    <div class="h-2 w-full bg-gray-100 rounded-full"></div>
                    <div class="h-2 w-5/6 bg-gray-100 rounded-full"></div>
                    <div class="h-2 w-4/6 bg-gray-100 rounded-full"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- 装饰球 -->
          <div class="absolute -top-10 -right-10 w-24 h-24 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 shadow-clay-btn animate-breathe flex items-center justify-center text-white text-3xl font-bold">
            Hi!
          </div>
        </div>
      </div>
    </header>

    <!-- Features Section -->
    <section class="py-20 px-6 bg-[#f0f3f8]/50">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="font-heading font-extrabold text-4xl sm:text-5xl text-clay-foreground mb-4">
            核心功能 <span class="text-clay-accent">一应俱全</span>
          </h2>
          <p class="text-xl text-clay-muted max-w-2xl mx-auto">从权限管理到代码生成，从系统监控到消息通知，常见中后台需求全部内置。</p>
        </div>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div
            v-for="(feature, index) in features"
            :key="index"
            class="group relative rounded-[32px] bg-white/70 p-8 shadow-clay-card backdrop-blur-xl transition-all duration-500 hover:-translate-y-2 hover:shadow-clay-card-hover overflow-hidden border border-[#d1d9e6]/40"
          >
            <div class="absolute -bottom-8 -right-8 w-32 h-32 rounded-full opacity-20 blur-2xl transition-transform duration-500 group-hover:scale-150" :class="feature.bg"></div>
            <div class="relative z-10">
              <div class="w-16 h-16 rounded-2xl bg-white shadow-clay-btn flex items-center justify-center mb-6 transition-transform duration-300 group-hover:scale-110 group-hover:rotate-6">
                <ArtSvgIcon :icon="feature.icon" class="text-[28px]" :class="feature.color" />
              </div>
              <h3 class="font-heading font-bold text-2xl text-clay-foreground mb-3">{{ feature.title }}</h3>
              <p class="text-clay-muted leading-relaxed">{{ feature.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CRUD Code Gen Section（对齐 homesite） -->
    <section class="py-20 px-6 relative">
      <div class="max-w-7xl mx-auto rounded-[60px] bg-white/30 backdrop-blur-xl shadow-clay-deep p-8 md:p-16 overflow-hidden relative">
        <div class="grid lg:grid-cols-2 gap-12 items-center overflow-hidden">
          <div class="order-2 lg:order-1 min-w-0 overflow-hidden">
            <div class="relative rounded-[32px] bg-[#1e1e1e] p-6 shadow-none lg:shadow-2xl overflow-hidden transform rotate-1 transition-transform hover:rotate-0">
              <div class="flex gap-2 mb-4">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <pre class="font-mono text-sm text-gray-300 overflow-x-auto"><span class="text-purple-400">export</span> <span class="text-blue-400">default</span> <span class="text-yellow-300">defineTable</span>({
  <span class="text-blue-300">columns</span>: [
    { <span class="text-green-300">prop</span>: <span class="text-orange-300">'id'</span>, <span class="text-green-300">label</span>: <span class="text-orange-300">'ID'</span>, <span class="text-green-300">width</span>: <span class="text-orange-300">80</span> },
    { <span class="text-green-300">prop</span>: <span class="text-orange-300">'title'</span>, <span class="text-green-300">label</span>: <span class="text-orange-300">'标题'</span> },
    { <span class="text-green-300">prop</span>: <span class="text-orange-300">'image'</span>, <span class="text-green-300">label</span>: <span class="text-orange-300">'图片'</span>, <span class="text-green-300">render</span>: <span class="text-yellow-300">Image</span> },
    { <span class="text-green-300">prop</span>: <span class="text-orange-300">'status'</span>, <span class="text-green-300">label</span>: <span class="text-orange-300">'状态'</span>, <span class="text-green-300">type</span>: <span class="text-orange-300">'switch'</span> },
  ],
  <span class="text-gray-500">// 自动生成 CRUD API 和 界面</span>
})</pre>
            </div>
          </div>
          <div class="order-1 lg:order-2 min-w-0">
            <div class="inline-block px-4 py-1 rounded-full bg-clay-accent/10 text-clay-accent font-bold mb-4">代码生成</div>
            <h2 class="font-heading font-extrabold text-4xl md:text-5xl text-clay-foreground mb-6">可视化 CRUD 代码生成</h2>
            <p class="text-lg text-clay-muted mb-8">在界面上配置数据表字段和表单类型，一键生成前后端完整的增删改查代码，包含控制器、模型、Vue 页面和 API 文件，无需手写重复代码。</p>
            <ul class="space-y-4">
              <li v-for="item in ['自动生成 Go 控制器与数据模型', '自动生成 Vue3 列表与表单页面', '支持关联查询、权限控制、导出等配置']" :key="item" class="flex items-center gap-3">
                <div class="w-6 h-6 rounded-full bg-clay-success/20 text-clay-success flex items-center justify-center">
                  <ArtSvgIcon icon="ri:check-line" class="text-sm" />
                </div>
                <span class="font-medium text-clay-foreground">{{ item }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'FrontendHome' })

const features = [
  {
    title: '权限管理',
    desc: '基于 RBAC 的角色权限体系，菜单、按钮、接口三级精确控制，支持多角色叠加。',
    icon: 'ri:shield-keyhole-line',
    bg: 'bg-blue-400',
    color: 'text-blue-500'
  },
  {
    title: '代码生成',
    desc: '可视化配置表结构，一键生成前后端增删改查完整代码，告别重复劳动。',
    icon: 'ri:tools-line',
    bg: 'bg-cyan-400',
    color: 'text-cyan-500'
  },
  {
    title: 'Vue3 + GoFrame',
    desc: '前端 Vue3 + TypeScript，后端 GoFrame v2，前后端分离，类型安全，开发体验极佳。',
    icon: 'ri:terminal-box-line',
    bg: 'bg-indigo-400',
    color: 'text-indigo-500'
  },
  {
    title: '系统监控',
    desc: '实时查看服务器 CPU、内存、磁盘状态，接口性能分析，运行日志一目了然。',
    icon: 'ri:dashboard-line',
    bg: 'bg-emerald-400',
    color: 'text-emerald-500'
  },
  {
    title: '消息通知',
    desc: '内置站内通知系统，支持系统公告、业务消息推送，WebSocket 实时送达。',
    icon: 'ri:notification-3-line',
    bg: 'bg-orange-400',
    color: 'text-orange-500'
  },
  {
    title: '开箱即用',
    desc: '完整的登录鉴权、操作日志、定时任务、文件上传等基础模块，克隆即可启动。',
    icon: 'ri:rocket-2-line',
    bg: 'bg-purple-400',
    color: 'text-purple-500'
  }
]
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.bg-clay-accent { background-color: #5a8dee; }
.bg-clay-success { background-color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.perspective-1000 { perspective: 1000px; }
.rotate-y-12 { transform: rotateY(12deg); }
.hover\:rotate-y-0:hover { transform: rotateY(0deg); }

.clay-text-gradient {
  background: linear-gradient(135deg, #32325d 20%, #5a8dee 60%, #03c3ec 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
}
.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-card-hover {
  box-shadow: 20px 20px 40px rgba(165, 175, 190, 0.35), -12px -12px 30px rgba(255, 255, 255, 0.95),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-btn-hover {
  box-shadow: 16px 16px 32px rgba(90, 141, 238, 0.4), -10px -10px 20px rgba(255, 255, 255, 0.5),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-pressed {
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
}

@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-20px); } }
@keyframes breathe { 0%, 100% { transform: scale(1); } 50% { transform: scale(1.05); } }
.animate-float { animation: float 8s ease-in-out infinite; }
.animate-float-delayed { animation: float 8s ease-in-out 4s infinite; }
.animate-breathe { animation: breathe 6s ease-in-out infinite; }
</style>
