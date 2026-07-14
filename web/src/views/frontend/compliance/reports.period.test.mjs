import assert from 'node:assert/strict'
import fs from 'node:fs'
import test from 'node:test'

const source = fs.readFileSync(new URL('./reports.vue', import.meta.url), 'utf8')

test('default report period uses the local calendar month', () => {
  assert.doesNotMatch(source, /toISOString\(\)\.slice\(0,\s*7\)/)
  const match = source.match(
    /function currentLocalPeriod\(date = new Date\(\)\) \{([\s\S]*?)\n\s{2}\}/
  )
  assert.ok(match, 'currentLocalPeriod helper should be defined in reports.vue')

  const currentLocalPeriod = Function(
    'return function currentLocalPeriod(date = new Date()) {' + match[1] + '\n}'
  )()
  // 2026-06-30 16:00 UTC is 2026-07-01 00:00 in Asia/Shanghai.
  assert.equal(currentLocalPeriod(new Date('2026-06-30T16:00:00.000Z')), '2026-07')
  assert.match(source, /ref\(currentLocalPeriod\(\)\)/)
})
