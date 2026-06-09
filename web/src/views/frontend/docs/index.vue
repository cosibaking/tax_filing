<!-- 政策与合规科普 — 文档中心 -->
<template>
  <main class="docs-page">
    <div class="docs-page__layout">
      <!-- 左侧导航 -->
      <aside class="docs-sidebar">
        <div class="docs-panel docs-panel--sidebar">
          <button type="button" class="docs-search-trigger" @click="openSearch">
            <ArtSvgIcon icon="ri:search-line" />
            <span>搜索文档...</span>
            <kbd>Ctrl K</kbd>
          </button>

          <div v-for="group in categories" :key="group.id" class="docs-nav-group">
            <h3 class="docs-nav-group__title">{{ group.title }}</h3>
            <ul class="docs-nav-list">
              <li v-for="child in (group.children || [])" :key="child.id">
                <button
                  type="button"
                  class="docs-nav-item"
                  :class="getCategoryClass(child)"
                  @click="toggleCategory(child)"
                >
                  <ArtSvgIcon :icon="child.icon || 'ri:folder-line'" />
                  <span>{{ child.title }}</span>
                </button>
                <ul
                  v-if="expandedCategoryId === child.id && getDocsByCategory(child.id).length > 1"
                  class="docs-nav-sublist"
                >
                  <li v-for="doc in getDocsByCategory(child.id)" :key="doc.slug">
                    <button
                      type="button"
                      class="docs-nav-item docs-nav-item--sub"
                      :class="{ 'is-active': activeDocSlug === doc.slug }"
                      @click="loadDoc(doc.slug)"
                    >
                      <ArtSvgIcon icon="ri:article-line" />
                      <span>{{ doc.title }}</span>
                    </button>
                  </li>
                </ul>
              </li>
            </ul>
          </div>

          <div v-if="loadingCategories" class="docs-empty-hint">加载中...</div>
        </div>
      </aside>

      <!-- 正文 -->
      <article class="docs-main">
        <div v-if="loadingDoc" class="docs-panel docs-panel--loading">
          <ArtSvgIcon icon="ri:loader-4-line" class="docs-spin" />
          <span>加载中...</span>
        </div>

        <div v-else-if="!currentDoc" class="docs-panel docs-panel--empty">
          <ArtSvgIcon icon="ri:file-text-line" class="docs-empty-icon" />
          <p>请从左侧选择一篇文档</p>
        </div>

        <div v-else class="docs-panel">
          <header class="docs-article__head">
            <h1>{{ currentDoc.title }}</h1>
            <div class="docs-article__meta">
              <span v-if="currentDoc.author">
                <ArtSvgIcon icon="ri:user-3-line" /> {{ currentDoc.author }}
              </span>
              <span v-if="currentDoc.updatedAt">
                <ArtSvgIcon icon="ri:time-line" /> {{ formatDate(currentDoc.updatedAt) }}
              </span>
              <span v-if="currentDoc.categoryName">
                <ArtSvgIcon icon="ri:bookmark-line" /> {{ currentDoc.categoryName }}
              </span>
            </div>
            <p v-if="currentDoc.summary" class="docs-article__summary">{{ currentDoc.summary }}</p>
          </header>

          <div class="docs-markdown">
            <MdPreview
              :model-value="currentDoc.content || ''"
              :editor-id="previewId"
              preview-theme="github"
              :no-mermaid="true"
              :no-katex="true"
            />
          </div>
        </div>
      </article>

      <!-- 右侧目录 -->
      <aside v-if="currentDoc" class="docs-toc-sidebar">
        <div class="docs-panel docs-panel--toc">
          <h4 class="docs-toc__title">本章目录</h4>
          <MdCatalog
            :key="activeDocSlug"
            :editor-id="previewId"
            :scroll-element="scrollElement"
            class="docs-toc"
          />
        </div>
      </aside>
    </div>

    <!-- 搜索弹窗 -->
    <Teleport to="body">
      <Transition name="docs-search-fade">
        <div v-if="searchVisible" class="docs-search-overlay" @mousedown.self="closeSearch">
          <div class="docs-search-dialog">
            <div class="docs-search-dialog__head">
              <ArtSvgIcon icon="ri:search-line" />
              <input
                ref="searchInputRef"
                v-model="searchKeyword"
                type="text"
                placeholder="搜索文档标题和内容..."
                @keydown.down.prevent="moveActive(1)"
                @keydown.up.prevent="moveActive(-1)"
                @keydown.enter.prevent="selectActive"
                @keydown.esc.prevent="closeSearch"
              />
              <kbd>ESC</kbd>
            </div>
            <div ref="searchBodyRef" class="docs-search-dialog__body">
              <div v-if="!searchKeyword.trim()" class="docs-search-tip">
                输入关键词搜索政策与合规文档
              </div>
              <div v-else-if="searchLoading" class="docs-search-tip">搜索中...</div>
              <div v-else-if="searchResults.length === 0" class="docs-search-tip">
                未找到「{{ searchKeyword.trim() }}」相关文档
              </div>
              <ul v-else class="docs-search-results">
                <li
                  v-for="(item, idx) in searchResults"
                  :key="item.slug"
                  class="docs-search-result"
                  :class="{ 'is-active': activeIndex === idx }"
                  @mouseenter="activeIndex = idx"
                  @click="selectResult(item)"
                >
                  <div>
                    <p class="docs-search-result__title" v-html="highlightKeyword(item.title)" />
                    <p class="docs-search-result__cat">
                      {{ item.categoryName }}
                      <span v-if="item.matchType === 'content'"> · 内容匹配</span>
                    </p>
                  </div>
                  <ArtSvgIcon icon="ri:arrow-right-s-line" />
                </li>
              </ul>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </main>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdPreview, MdCatalog } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { fetchDocCategoryTree, fetchDocListByCategory, fetchDocDetailBySlug, fetchDocSearch } from '@/api/frontend/doc'
import {
  buildPolicyDocCategoryTree,
  getPolicyDocsByCategory,
  getPolicyDocBySlug,
  searchPolicyDocs,
  getFirstPolicyDocSlug,
  POLICY_DOC_CATEGORIES,
} from '@/data/frontend/policyDocs'

defineOptions({ name: 'FrontendDocs' })

const route = useRoute()
const router = useRouter()
const previewId = 'doc-preview'
const scrollElement = document.documentElement

const categories = ref<any[]>([])
const docsByCategory = ref<Record<string, any[]>>({})
const loadingCategories = ref(false)
const loadingDoc = ref(false)
const activeDocSlug = ref('')
const currentDoc = ref<any>(null)
const expandedCategoryId = ref<string | number>()

const searchVisible = ref(false)
const searchKeyword = ref('')
const searchResults = ref<any[]>([])
const searchLoading = ref(false)
const activeIndex = ref(0)
const searchInputRef = ref<HTMLInputElement>()
const searchBodyRef = ref<HTMLDivElement>()
let searchTimer: ReturnType<typeof setTimeout> | null = null

function getDocsByCategory(categoryId: string | number) {
  return docsByCategory.value[String(categoryId)] || []
}

function getCategoryClass(child: any) {
  const docs = getDocsByCategory(child.id)
  const isSingle = docs.length <= 1
  const isExpanded = expandedCategoryId.value === child.id
  if (isSingle && isExpanded && docs[0]?.slug && activeDocSlug.value === docs[0].slug) {
    return 'is-active'
  }
  if (isExpanded) return 'is-expanded'
  return ''
}

function toggleCategory(child: any) {
  const docs = getDocsByCategory(child.id)
  if (docs.length <= 1) {
    expandedCategoryId.value = child.id
    if (docs.length === 1 && docs[0].slug) loadDoc(docs[0].slug)
    return
  }
  expandedCategoryId.value = expandedCategoryId.value === child.id ? undefined : child.id
}

async function loadDoc(slug: string) {
  if (!slug || (activeDocSlug.value === slug && currentDoc.value)) return
  activeDocSlug.value = slug
  loadingDoc.value = true
  router.replace({ path: '/docs', query: { slug } })

  try {
    const builtIn = getPolicyDocBySlug(slug)
    if (builtIn) {
      currentDoc.value = builtIn
      window.scrollTo({ top: 0, behavior: 'smooth' })
      return
    }
    const detail = await fetchDocDetailBySlug(slug)
    currentDoc.value = detail
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } catch (e) {
    console.error('加载文档失败', e)
    currentDoc.value = null
  } finally {
    loadingDoc.value = false
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function initBuiltInDocs() {
  categories.value = buildPolicyDocCategoryTree()
  for (const cat of POLICY_DOC_CATEGORIES) {
    docsByCategory.value[cat.id] = getPolicyDocsByCategory(cat.id)
  }
}

async function loadCmsDocs() {
  try {
    const tree = await fetchDocCategoryTree()
    if (!tree?.length) return
    const cmsDocs: Record<string, any[]> = {}
    const childIds: number[] = []
    for (const group of tree) {
      for (const child of group.children || []) {
        if (child.id) childIds.push(child.id)
      }
    }
    const results = await Promise.all(childIds.map((id) => fetchDocListByCategory(id)))
    childIds.forEach((id, idx) => {
      cmsDocs[String(id)] = results[idx] || []
    })
    categories.value.push(...tree)
    docsByCategory.value = { ...docsByCategory.value, ...cmsDocs }
  } catch {
    /* CMS 文档可选 */
  }
}

async function loadCategories() {
  loadingCategories.value = true
  initBuiltInDocs()
  await loadCmsDocs()

  const querySlug = typeof route.query.slug === 'string' ? route.query.slug : ''
  const firstSlug = querySlug || getFirstPolicyDocSlug()
  if (firstSlug) {
    const article = getPolicyDocBySlug(firstSlug)
    if (article) {
      expandedCategoryId.value = POLICY_DOC_CATEGORIES.find((c) =>
        getPolicyDocsByCategory(c.id).some((d) => d.slug === firstSlug)
      )?.id
    }
    await loadDoc(firstSlug)
  }
  loadingCategories.value = false
}

function openSearch() {
  searchVisible.value = true
  searchKeyword.value = ''
  searchResults.value = []
  activeIndex.value = 0
  nextTick(() => searchInputRef.value?.focus())
}

function closeSearch() {
  searchVisible.value = false
  searchKeyword.value = ''
}

function moveActive(delta: number) {
  if (!searchResults.value.length) return
  activeIndex.value = (activeIndex.value + delta + searchResults.value.length) % searchResults.value.length
}

function selectActive() {
  const item = searchResults.value[activeIndex.value]
  if (item) selectResult(item)
}

function selectResult(item: any) {
  closeSearch()
  loadDoc(item.slug)
}

function highlightKeyword(text: string) {
  const kw = searchKeyword.value.trim()
  if (!kw) return text
  const escaped = kw.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return text.replace(new RegExp(`(${escaped})`, 'gi'), '<mark>$1</mark>')
}

watch(searchKeyword, (kw) => {
  if (searchTimer) clearTimeout(searchTimer)
  const trimmed = kw.trim()
  if (!trimmed) {
    searchResults.value = []
    searchLoading.value = false
    return
  }
  searchLoading.value = true
  searchTimer = setTimeout(async () => {
    try {
      const builtIn = searchPolicyDocs(trimmed)
      let apiResults: any[] = []
      try {
        apiResults = await fetchDocSearch(trimmed)
      } catch { /* ignore */ }
      const seen = new Set(builtIn.map((r) => r.slug))
      searchResults.value = [
        ...builtIn,
        ...apiResults.filter((r) => r.slug && !seen.has(r.slug)),
      ]
      activeIndex.value = 0
    } finally {
      searchLoading.value = false
    }
  }, 250)
})

function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchVisible.value ? closeSearch() : openSearch()
  }
}

watch(() => route.query.slug, (slug) => {
  if (typeof slug === 'string' && slug && slug !== activeDocSlug.value) {
    loadDoc(slug)
  }
})

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeydown)
  loadCategories()
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
  if (searchTimer) clearTimeout(searchTimer)
})
</script>

<style lang="scss" scoped>
.docs-page {
  padding: 24px;
  max-width: 1280px;
  margin: 0 auto;
}

.docs-page__layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 20px;
  align-items: start;

  @media (min-width: 1280px) {
    grid-template-columns: 240px 1fr 200px;
  }
}

.docs-panel {
  padding: 20px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.docs-sidebar {
  position: sticky;
  top: 88px;
  max-height: calc(100vh - 100px);
  overflow-y: auto;
}

.docs-panel--sidebar {
  padding: 16px;
}

.docs-search-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  margin-bottom: 16px;
  padding: 8px 12px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
  font-size: 12px;
  color: #64748b;
  cursor: pointer;

  span { flex: 1; text-align: left; }

  kbd {
    padding: 2px 6px;
    border-radius: 4px;
    border: 1px solid #d8dee9;
    background: #fff;
    font-size: 10px;
  }

  &:hover {
    border-color: #bfdbfe;
    color: #2563eb;
  }
}

.docs-nav-group__title {
  margin: 0 0 8px;
  padding: 0 8px;
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.docs-nav-list,
.docs-nav-sublist {
  margin: 0 0 12px;
  padding: 0;
  list-style: none;
}

.docs-nav-sublist {
  margin-left: 12px;
}

.docs-nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 10px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: #f8fafc;
    color: #2563eb;
  }

  &.is-expanded {
    background: #eff6ff;
    color: #2563eb;
  }

  &.is-active {
    background: #2563eb;
    color: #fff;
  }

  &--sub {
    font-size: 12px;
    font-weight: 500;
  }
}

.docs-panel--loading,
.docs-panel--empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  gap: 12px;
  color: #64748b;
  font-size: 14px;
}

.docs-empty-icon {
  font-size: 40px;
  color: #94a3b8;
}

.docs-spin {
  font-size: 28px;
  color: #2563eb;
  animation: spin 1s linear infinite;
}

.docs-article__head {
  padding-bottom: 20px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e8edf3;

  h1 {
    margin: 0 0 12px;
    font-size: 24px;
    font-weight: 800;
    color: #1a1f36;
    line-height: 1.35;
  }
}

.docs-article__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;

  span {
    display: inline-flex;
    align-items: center;
    gap: 4px;
  }
}

.docs-article__summary {
  margin: 14px 0 0;
  font-size: 14px;
  line-height: 1.6;
  color: #64748b;
}

.docs-toc-sidebar {
  display: none;

  @media (min-width: 1280px) {
    display: block;
    position: sticky;
    top: 88px;
    max-height: calc(100vh - 100px);
    overflow-y: auto;
  }
}

.docs-panel--toc {
  padding: 16px;
}

.docs-toc__title {
  margin: 0 0 12px;
  padding: 0 8px;
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
}

:deep(.docs-markdown) {
  .md-editor-preview-wrapper { padding: 0; }
  .md-editor-preview {
    font-size: 14px;
    line-height: 1.8;
    color: #475569;

    h2 {
      margin: 28px 0 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid #f1f5f9;
      font-size: 17px;
      font-weight: 700;
      color: #1a1f36;
    }
    h3 {
      margin: 20px 0 8px;
      font-size: 15px;
      font-weight: 700;
      color: #334155;
    }
    p { margin: 0 0 12px; }
    strong { color: #1a1f36; }
    ul, ol { margin: 0 0 12px; padding-left: 1.4em; }
    li { margin-bottom: 6px; }
    table {
      width: 100%;
      margin: 12px 0;
      border-collapse: collapse;
      font-size: 13px;
    }
    th, td {
      padding: 10px 12px;
      border: 1px solid #e8edf3;
      text-align: left;
    }
    th { background: #f8fafc; font-weight: 700; }
    blockquote {
      margin: 12px 0;
      padding: 12px 16px;
      border-left: 4px solid #2563eb;
      background: #eff6ff;
      border-radius: 0 8px 8px 0;
      color: #475569;
    }
    a { color: #2563eb; }
  }
}

:deep(.docs-toc) {
  .md-editor-catalog-link {
    font-size: 12px;
    font-weight: 600;
    color: #94a3b8;
    padding: 4px 8px;
    border-radius: 6px;
    &.md-editor-catalog-active {
      color: #2563eb;
      background: #eff6ff;
    }
  }
}

.docs-search-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
  background: rgba(15, 23, 42, 0.35);
  backdrop-filter: blur(4px);
}

.docs-search-dialog {
  width: 100%;
  max-width: 560px;
  margin: 0 16px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e8edf3;
  box-shadow: 0 20px 50px rgba(15, 23, 42, 0.15);
  overflow: hidden;
}

.docs-search-dialog__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid #e8edf3;

  input {
    flex: 1;
    border: none;
    outline: none;
    font-size: 14px;
    color: #1a1f36;
  }

  kbd {
    padding: 2px 6px;
    border-radius: 4px;
    border: 1px solid #d8dee9;
    font-size: 10px;
    color: #94a3b8;
  }
}

.docs-search-dialog__body {
  max-height: 400px;
  overflow-y: auto;
}

.docs-search-tip {
  padding: 40px 20px;
  text-align: center;
  font-size: 13px;
  color: #94a3b8;
}

.docs-search-results {
  margin: 0;
  padding: 8px;
  list-style: none;
}

.docs-search-result {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;

  &:hover, &.is-active {
    background: #eff6ff;
  }

  &__title {
    margin: 0;
    font-size: 14px;
    font-weight: 700;
    color: #1a1f36;
  }

  &__cat {
    margin: 4px 0 0;
    font-size: 11px;
    color: #94a3b8;
  }

  :deep(mark) {
    background: rgba(37, 99, 235, 0.15);
    color: #2563eb;
    border-radius: 2px;
  }
}

.docs-search-fade-enter-active,
.docs-search-fade-leave-active {
  transition: opacity 0.15s ease;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 960px) {
  .docs-page__layout {
    grid-template-columns: 1fr;
  }

  .docs-sidebar {
    position: static;
    max-height: none;
  }
}
</style>
