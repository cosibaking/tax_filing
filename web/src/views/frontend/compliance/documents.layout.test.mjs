import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const source = readFileSync(new URL('./documents.vue', import.meta.url), 'utf8')

test('documents upload row constrains all three desktop columns', () => {
  assert.match(source, /class="period-field"/)
  assert.match(source, /class="upload-field"/)
  assert.match(source, /grid-template-columns:\s*180px minmax\(0, 1fr\) auto/)
  assert.match(source, /\.upload-field\s*\{[^}]*min-width:\s*0/s)
  assert.match(source, /\.period-field\s*\{[^}]*--el-date-editor-width:\s*100%/s)
})

test('documents upload row becomes full-width single column on mobile', () => {
  assert.match(source, /@media \(max-width: 760px\)/)
  assert.match(source, /\.period-field\s*[^}]*width:\s*100%/s)
})
