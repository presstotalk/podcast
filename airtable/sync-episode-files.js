const table = base.getTable('Episodes')
const titleField = table.getField('Title')
const slugField = table.getField('Slug')

const record = await input.recordAsync('Choose a record', table)

if (record) {
  const title = record.getCellValue(titleField)
  const slug = record.getCellValue(slugField)
  await remoteFetchAsync(`https://presstotalk-syncer.poying-me.workers.dev/${encodeURIComponent(title)}/cover.jpg?dest=${encodeURIComponent(`${encodeURIComponent(slug)}/cover.jpg`)}`, { method: 'POST' })
  await remoteFetchAsync(`https://presstotalk-syncer.poying-me.workers.dev/${encodeURIComponent(title)}/audio.mp3?dest=${encodeURIComponent(`${encodeURIComponent(slug)}/audio.mp3`)}`, { method: 'POST' })
  output.markdown(`success`)
}
