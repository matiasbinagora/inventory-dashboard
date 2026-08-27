'use client'

import { useEffect, useMemo, useState } from 'react'
import { countUpValue, filterProjects } from './catalog'
import type { Project } from './types'

const icons = { dashboard: '⌂', catalog: '▦', admin: '⚙' }

function CountUp({ value }: { value: number }) {
  const [display, setDisplay] = useState(0)
  useEffect(() => {
    const start = performance.now()
    let frame = 0
    const tick = (now: number) => {
      setDisplay(countUpValue(value, now - start))
      if (now - start < 700) frame = requestAnimationFrame(tick)
      else setDisplay(value)
    }
    frame = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(frame)
  }, [value])
  return <span aria-label={`${value}`}>{display}</span>
}

function Sidebar({ collapsed, onToggle }: { collapsed: boolean; onToggle: () => void }) {
  return <aside className={`sidebar ${collapsed ? 'sidebar-collapsed' : ''}`} aria-label="Primary navigation">
    <div className="brand"><span className="brand-mark">◈</span><span className="sidebar-label">INVENTORY</span></div>
    <nav><a className="nav-item active" href="#dashboard"><span aria-hidden="true">{icons.dashboard}</span><span className="sidebar-label">Dashboard</span></a><a className="nav-item" href="#catalog"><span aria-hidden="true">{icons.catalog}</span><span className="sidebar-label">Catalog</span></a><a className="nav-item" href="#administration" aria-disabled="true"><span aria-hidden="true">{icons.admin}</span><span className="sidebar-label">Administration</span></a></nav>
    <button className="collapse-button" onClick={onToggle} aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'} aria-expanded={!collapsed}>{collapsed ? '→' : '←'}<span className="sidebar-label">Collapse</span></button>
  </aside>
}

function Metric({ label, value, tone }: { label: string; value: number; tone: string }) {
  return <article className={`metric metric-${tone}`}><p>{label}</p><strong><CountUp value={value} /></strong><span className="metric-line" /></article>
}

function Thumbnail({ project }: { project: Project }) {
  const thumbnail = project.media.find((media) => media.role === 'thumbnail')
  return thumbnail ? <img src={thumbnail.source} alt={thumbnail.alt_text || `${project.name} thumbnail`} /> : <div className="thumbnail-placeholder" aria-hidden="true">{project.name.slice(0, 1)}</div>
}

function ProjectCard({ project }: { project: Project }) {
  return <a className="project-card" href={`/projects/${encodeURIComponent(project.id)}`}><Thumbnail project={project} /><div className="card-content"><span className="eyebrow">PROJECT / {project.id.slice(0, 8)}</span><h3>{project.name}</h3><p>{project.description || 'Curated project metadata'}</p><div className="tag-row">{project.technologies.slice(0, 3).map((technology) => <span key={technology} className="tag">{technology}</span>)}</div></div><span className="card-arrow" aria-hidden="true">↗</span></a>
}

function CatalogTable({ projects }: { projects: Project[] }) {
  return <div className="table-wrap"><table><thead><tr><th>Project</th><th>Description</th><th>Technologies</th><th>Links</th></tr></thead><tbody>{projects.map((project) => <tr key={project.id}><td><a href={`/projects/${encodeURIComponent(project.id)}`} className="table-project"><Thumbnail project={project} /><strong>{project.name}</strong></a></td><td>{project.description || '—'}</td><td>{project.technologies.join(', ') || '—'}</td><td>{project.links.length || '—'}</td></tr>)}</tbody></table></div>
}

export function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [collapsed, setCollapsed] = useState(false)
  const [view, setView] = useState<'cards' | 'table'>('cards')
  const [search, setSearch] = useState('')
  const [technology, setTechnology] = useState('')
  const [platform, setPlatform] = useState('')

  const loadProjects = async () => {
    setLoading(true); setError(false)
    try { const response = await fetch('/api/projects', { cache: 'no-store' }); if (!response.ok) throw new Error('API error'); setProjects(await response.json()) } catch { setError(true) } finally { setLoading(false) }
  }
  useEffect(() => { void loadProjects() }, [])
  const technologies = useMemo(() => [...new Set(projects.flatMap((project) => project.technologies))].sort(), [projects])
  const platforms = useMemo(() => [...new Set(projects.map((project) => project.agentic_platform).filter(Boolean))].sort(), [projects])
  const filtered = useMemo(() => filterProjects(projects, search, technology, platform), [projects, search, technology, platform])
  const mediaProjects = projects.filter((project) => project.media.length > 0).length
  const linkProjects = projects.filter((project) => project.links.length > 0).length

  return <div className="app-shell"><Sidebar collapsed={collapsed} onToggle={() => setCollapsed(!collapsed)} /><main className="main-content" id="dashboard"><header className="topbar"><div><span className="eyebrow">LOCAL / EDITORIAL CONTROL ROOM</span><h1>Project inventory</h1><p>Signal from the workbench, curated for fast discovery.</p></div><div className="status-pill"><span />LOCAL API</div></header>
    <section className="metrics" aria-label="Inventory metrics"><Metric label="Projects cataloged" value={projects.length} tone="violet" /><Metric label="With curated media" value={mediaProjects} tone="orange" /><Metric label="Public link sets" value={linkProjects} tone="cyan" /><Metric label="Visible results" value={filtered.length} tone="lime" /></section>
    <section className="catalog-section" id="catalog"><div className="section-heading"><div><span className="eyebrow">02 / CATALOG</span><h2>Editorial index</h2></div><div className="view-toggle" role="group" aria-label="Catalog view"><button className={view === 'cards' ? 'selected' : ''} onClick={() => setView('cards')} aria-pressed={view === 'cards'}>▦ Cards</button><button className={view === 'table' ? 'selected' : ''} onClick={() => setView('table')} aria-pressed={view === 'table'}>☷ Table</button></div></div>
      <div className="filters"><label className="search-field"><span aria-hidden="true">⌕</span><input aria-label="Search projects" placeholder="Search projects, descriptions, technologies…" value={search} onChange={(event) => setSearch(event.target.value)} /></label><label>Technology<select aria-label="Filter by technology" value={technology} onChange={(event) => setTechnology(event.target.value)}><option value="">All technologies</option>{technologies.map((item) => <option key={item}>{item}</option>)}</select></label><label>Platform<select aria-label="Filter by platform" value={platform} onChange={(event) => setPlatform(event.target.value)}><option value="">All platforms</option>{platforms.map((item) => <option key={item}>{item}</option>)}</select></label></div>
      {loading ? <div className="loading-state" role="status">Loading catalog…</div> : error ? <div className="error-state" role="alert"><h3>Could not reach the local API</h3><p>Start the Go service on 127.0.0.1:8080 and try again.</p><button onClick={() => void loadProjects()}>Retry</button></div> : filtered.length === 0 ? <div className="empty-state"><span>◌</span><h3>{projects.length ? 'No projects match these filters' : 'Your catalog is ready for its first project'}</h3><p>{projects.length ? 'Try a broader search or clear the filters.' : 'Curated metadata from the local API will appear here.'}</p></div> : view === 'cards' ? <div className="project-grid">{filtered.map((project) => <ProjectCard key={project.id} project={project} />)}</div> : <CatalogTable projects={filtered} />}
    </section></main></div>
}
