export default async function ProjectPreview({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  return <main className="page-shell"><a href="/">← Back to catalog</a><div className="detail-preview"><span className="eyebrow">PROJECT / {id}</span><h1>Project detail</h1><p>The full editorial project view is prepared for TASK-006.</p></div></main>
}
