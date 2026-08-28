'use client'

import { useEffect, useMemo, useState } from 'react'
import type { Media, Project } from '../dashboard/types'
import { fullSizeSource, imageMedia, localMediaURL, parseMetadataEntries, sortMilestones } from './detail'

function ProjectDetailError() {
  return <main className="page-shell"><a href="/">← Back to catalog</a><div className="detail-state" role="alert"><h1>Project unavailable</h1><p>The local API could not load this project.</p></div></main>
}

function LinkList({ project }: { project: Project }) {
  const links = project.links.filter((link) => link.kind === 'github' || link.kind === 'trello')
  if (!links.length) return null
  return <section className="detail-section public-links-section" aria-labelledby="links-heading"><span className="eyebrow">PUBLIC REFERENCES</span><h2 id="links-heading">Continue exploring</h2><div className="public-links">{links.map((link) => <a key={link.id ?? link.url} href={link.url} target="_blank" rel="noreferrer" className="public-link"><span>{link.label || (link.kind === 'github' ? 'GitHub repository' : 'Trello board')}</span><span aria-hidden="true">↗</span></a>)}</div></section>
}

function Gallery({ project, onSelect }: { project: Project; onSelect: (media: Media) => void }) {
  const images = imageMedia(project)
  if (!images.length) return null
  return <section className="detail-section" aria-labelledby="gallery-heading"><span className="eyebrow">CURATED MEDIA</span><h2 id="gallery-heading">Graphify Report</h2><div className="detail-gallery">{images.map((media) => <button type="button" className="gallery-item" key={media.id ?? media.source} onClick={() => onSelect(media)}><img src={localMediaURL(media.source)} alt={media.alt_text || `${project.name} ${media.role}`} /><span>{media.caption || (media.role === 'thumbnail' ? 'Open original' : media.role)}</span></button>)}</div></section>
}

function VideoReferences({ project }: { project: Project }) {
  const videos = project.media.filter((media) => media.role === 'video')
  if (!videos.length) return null
  return <section className="detail-section" aria-labelledby="video-heading"><span className="eyebrow">DEMO</span><h2 id="video-heading">Demo Video</h2><div className="video-list">{videos.map((video) => /^https:\/\//i.test(video.source) ? <a className="video-reference" key={video.id ?? video.source} href={video.source} target="_blank" rel="noreferrer"><span>{video.caption || 'Open public video'}</span><span aria-hidden="true">↗</span></a> : <video className="local-video" key={video.id ?? video.source} controls preload="metadata" aria-label={video.caption || `${project.name} demo`}><source src={localMediaURL(video.source)} /><track kind="captions" /></video>)}</div></section>
}

function AgenticPlatform({ value }: { value: string }) {
  return <ul className="metadata-list agentic-platform-list">{parseMetadataEntries(value).map((entry, index) => <li key={`${entry.primary}-${index}`}><span className="metadata-primary">{entry.primary}</span>{entry.urls.map((url) => <a className="metadata-link" key={url} href={url} target="_blank" rel="noreferrer">{url}</a>)}</li>)}</ul>
}

function MetadataList({ value, className = '' }: { value: string | string[]; className?: string }) {
  return <ul className={`metadata-list ${className}`.trim()}>{parseMetadataEntries(value).map((entry, index) => <li key={`${entry.primary}-${index}`}><span className="metadata-primary">{entry.primary}</span>{entry.urls.map((url) => <a className="metadata-link" key={url} href={url} target="_blank" rel="noreferrer">{url}</a>)}</li>)}</ul>
}

export function ProjectDetail({ projectID }: { projectID: string }) {
  const [project, setProject] = useState<Project | null>(null)
  const [selected, setSelected] = useState<Media | null>(null)
  const [error, setError] = useState(false)

  useEffect(() => {
    let active = true
    fetch(`/api/projects/${encodeURIComponent(projectID)}`, { cache: 'no-store' }).then((response) => {
      if (!response.ok) throw new Error('project unavailable')
      return response.json() as Promise<Project>
    }).then((value) => { if (active) setProject(value) }).catch(() => { if (active) setError(true) })
    return () => { active = false }
  }, [projectID])

  const milestones = useMemo(() => project ? sortMilestones(project.milestones) : [], [project])
  if (error) return <ProjectDetailError />
  if (!project) return <main className="page-shell"><div className="loading-state" role="status">Loading project detail…</div></main>

  const title = project.name.trim() || 'Untitled project'
  return <main className="detail-page"><a href="/" className="back-link">← Back to catalog</a><header className="detail-hero"><span className="eyebrow">PROJECT / {project.id.slice(0, 8)}</span><h1>{title}</h1>{project.description && <p>{project.description}</p>}</header><div className="detail-layout"><div className="detail-primary"><Gallery project={project} onSelect={setSelected} /><VideoReferences project={project} />{milestones.length > 0 && <section className="detail-section" aria-labelledby="history-heading"><span className="eyebrow">MANUAL HISTORY</span><h2 id="history-heading">The work, over time</h2><ol className="timeline">{milestones.map((milestone) => <li key={milestone.id ?? `${milestone.date}-${milestone.title}`}><time dateTime={milestone.date}>{milestone.date}</time><div><h3>{milestone.title}</h3><p>{milestone.description}</p></div></li>)}</ol></section>}</div><aside className="detail-aside"><section className="detail-section metadata-section"><span className="eyebrow">REGISTERED FIELDS</span>{project.technologies.length > 0 && <div className="field-group"><h2>Technologies</h2><MetadataList value={project.technologies} className="technology-list" /></div>}{project.agentic_platform && <div className="field-group agentic-platform"><h2>Agentic platform</h2><AgenticPlatform value={project.agentic_platform} /></div>}</section><LinkList project={project} /></aside></div>{selected && <div className="lightbox" role="dialog" aria-modal="true" aria-label="Expanded project image" onClick={() => setSelected(null)}><div className="lightbox-content" onClick={(event) => event.stopPropagation()}><button type="button" className="lightbox-close" aria-label="Close image" onClick={() => setSelected(null)}>×</button><img src={localMediaURL(fullSizeSource(selected, project.media))} alt={selected.alt_text || `${project.name} expanded`} />{selected.caption && <p>{selected.caption}</p>}</div></div>}</main>
}
