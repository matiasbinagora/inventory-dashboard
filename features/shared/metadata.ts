export type MetadataEntry = { primary: string; urls: string[] }

const URL_PATTERN = /https:\/\/[^\s,;]+/gi

function cleanURL(value: string) {
  return value.replace(/[.)\]}]+$/g, '')
}

function parseMetadataSegment(value: string): MetadataEntry[] {
  const matches = [...value.matchAll(URL_PATTERN)]
  if (!matches.length) return [{ primary: value.trim(), urls: [] }]
  const entries = matches.map((match, index) => {
    const previousEnd = index === 0 ? 0 : (matches[index - 1].index || 0) + matches[index - 1][0].length
    const primary = value.slice(previousEnd, match.index || 0).replace(/^[\s:—–,;-]+|[\s:—–,;-]+$/g, '').trim()
    return { primary: primary || (index === 0 ? value.trim() : 'Reference'), urls: [cleanURL(match[0])] }
  })
  const last = matches[matches.length - 1]
  const trailing = value.slice((last.index || 0) + last[0].length).replace(/^[\s:—–,;-]+|[\s:—–,;-]+$/g, '').trim()
  return trailing ? [...entries, { primary: trailing, urls: [] }] : entries
}

export function parseMetadataEntries(value: string | string[]) {
  const values = Array.isArray(value) ? value : [value]
  return values.flatMap((item) => item.split(/\n|;\s*/).map((segment) => segment.trim()).filter(Boolean).flatMap(parseMetadataSegment))
}
