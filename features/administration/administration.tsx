'use client'

import { useEffect, useState } from 'react'
import type { Media, Milestone, Project } from '../dashboard/types'
import { addMedia, addMilestone, getProject, listProjects, removeMedia, removeMilestone, saveProject, updateMilestone } from './api'
import { thumbnailPayload } from './form'

type ProjectForm = Omit<Project, 'id'>
const emptyProject: ProjectForm = { name: '', description: '', agentic_platform: '', technologies: [], links: [], media: [], milestones: [] }
const emptyMilestone = { date: '', title: '', description: '', media_ids: [] as string[] }
const roles: Media['role'][] = ['original', 'screenshot', 'video']

export function Administration() {
  const [projects, setProjects] = useState<Project[]>([])
  const [project, setProject] = useState<Project | null>(null)
  const [form, setForm] = useState<ProjectForm>(emptyProject)
  const [media, setMedia] = useState({ role: 'screenshot' as Media['role'], source: '', alt_text: '', caption: '' })
  const [milestone, setMilestone] = useState(emptyMilestone)
  const [editingMilestone, setEditingMilestone] = useState<string | null>(null)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  const refresh = async (id?: string) => {
    const all = await listProjects(); setProjects(all)
    if (id) { const loaded = await getProject(id); setProject(loaded); setForm({ name: loaded.name, description: loaded.description || '', agentic_platform: loaded.agentic_platform || '', technologies: loaded.technologies, links: loaded.links, media: loaded.media, milestones: loaded.milestones }) }
  }
  useEffect(() => { void refresh().catch((cause: Error) => setError(cause.message)) }, [])
  const selectProject = (id: string) => { if (!id) { setProject(null); setForm(emptyProject); return }; void refresh(id).catch((cause: Error) => setError(cause.message)) }
  const setField = (field: keyof typeof emptyProject, value: string) => setForm((current) => ({ ...current, [field]: value }))
  const save = async (event: React.FormEvent) => {
    event.preventDefault(); setMessage(''); setError('')
    try { const wasEditing = Boolean(project); const saved = await saveProject(form, project?.id); setProject(saved); await refresh(saved.id); setMessage(wasEditing ? 'Project updated.' : 'Project created.'); window.history.replaceState(null, '', `/admin?project=${saved.id}`) }
    catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not save project.') }
  }
  const addCuratedMedia = async (event: React.FormEvent) => {
    event.preventDefault(); if (!project || !media.source.trim()) return
    try { await addMedia(project.id, { ...media, curated: true }); await refresh(project.id); setMedia({ ...media, source: '' }); setMessage('Media association added.') } catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not add media.') }
  }
  const chooseThumbnail = async (item: Media) => {
    if (!project || !window.confirm('Use this curated asset as the main thumbnail? The original association will remain.')) return
    try { const current = project.media.find((candidate) => candidate.role === 'thumbnail'); if (current?.id) await removeMedia(project.id, current.id); await addMedia(project.id, thumbnailPayload(item)); await refresh(project.id); setMessage('Main thumbnail updated; original media preserved.') } catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not choose thumbnail.') }
  }
  const removeCuratedMedia = async (item: Media) => { if (!project || !item.id || !window.confirm('Remove only this media association?')) return; try { await removeMedia(project.id, item.id); await refresh(project.id); setMessage('Media association removed.') } catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not remove media.') } }
  const saveMilestone = async (event: React.FormEvent) => {
    event.preventDefault(); if (!project) return
    try { if (editingMilestone) await updateMilestone(project.id, editingMilestone, milestone); else await addMilestone(project.id, milestone); await refresh(project.id); setMilestone(emptyMilestone); setEditingMilestone(null); setMessage(editingMilestone ? 'Milestone updated.' : 'Milestone added.') } catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not save milestone.') }
  }
  const deleteHistory = async (item: Milestone) => { if (!project || !item.id || !window.confirm(`Delete milestone “${item.title}”?`)) return; try { await removeMilestone(project.id, item.id); await refresh(project.id); setMessage('Milestone deleted.') } catch (cause) { setError(cause instanceof Error ? cause.message : 'Could not delete milestone.') } }

  return <main className="admin-page"><header className="admin-header"><div className="admin-header-nav"><a href="/">← Back to dashboard</a><span className="eyebrow">LOCAL / ADMINISTRATION</span></div><h1>Inventory Administration</h1><p>Safe metadata, approved media, and manual history. No login required.</p></header>
    <section className="admin-toolbar"><label>Project<select aria-label="Select project" value={project?.id ?? ''} onChange={(event) => selectProject(event.target.value)}><option value="">New project</option>{projects.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}</select></label><span className="status-pill"><span />LOCAL API</span></section>
    {message && <p className="success-message" role="status">{message}</p>}{error && <p className="form-error" role="alert">{error}</p>}
    <form className="admin-card" onSubmit={save}><span className="eyebrow">01 / PROJECT METADATA</span><h2>{project ? 'Edit project' : 'Create project'}</h2><label>Name *<input required aria-label="Project name" value={form.name} onChange={(event) => setField('name', event.target.value)} /></label><label>Description<textarea aria-label="Project description" value={form.description} onChange={(event) => setField('description', event.target.value)} /></label><label>Agentic platform<input value={form.agentic_platform} onChange={(event) => setField('agentic_platform', event.target.value)} /></label><label>Technologies <input value={form.technologies.join(', ')} onChange={(event) => setForm({ ...form, technologies: event.target.value.split(',').map((value) => value.trim()).filter(Boolean) })} placeholder="Go, Next.js" /></label><button className="primary-button" type="submit">{project ? 'Save changes' : 'Create project'}</button></form>
    {project && <><section className="admin-card"><span className="eyebrow">02 / CURATED MEDIA</span><h2>Media associations</h2><p className="helper-text">Use managed <code>media/</code> paths or public HTTPS video URLs. Files are never copied from repositories.</p><form className="inline-form" onSubmit={addCuratedMedia}><label>Role<select value={media.role} onChange={(event) => setMedia({ ...media, role: event.target.value as Media['role'] })}>{roles.map((role) => <option key={role}>{role}</option>)}</select></label><label>Managed source *<input required value={media.source} onChange={(event) => setMedia({ ...media, source: event.target.value })} placeholder="media/project-shot.png" /></label><label>Alt text<input value={media.alt_text} onChange={(event) => setMedia({ ...media, alt_text: event.target.value })} /></label><label>Caption<input value={media.caption} onChange={(event) => setMedia({ ...media, caption: event.target.value })} /></label><button className="primary-button" type="submit">Add media</button></form><ul className="admin-list">{project.media.map((item) => <li key={item.id ?? item.source}><span><strong>{item.role}</strong> {item.source}{item.role === 'original' && <em> preserved original</em>}</span><span className="list-actions">{item.role !== 'thumbnail' && <button type="button" onClick={() => void chooseThumbnail(item)}>Use thumbnail</button>}<button type="button" onClick={() => void removeCuratedMedia(item)}>Remove association</button></span></li>)}</ul></section><section className="admin-card"><span className="eyebrow">03 / MANUAL HISTORY</span><h2>Milestones</h2><form className="milestone-form" onSubmit={saveMilestone}><label>Date *<input required type="date" value={milestone.date} onChange={(event) => setMilestone({ ...milestone, date: event.target.value })} /></label><label>Title *<input required value={milestone.title} onChange={(event) => setMilestone({ ...milestone, title: event.target.value })} /></label><label>Description *<textarea required value={milestone.description} onChange={(event) => setMilestone({ ...milestone, description: event.target.value })} /></label><button className="primary-button" type="submit">{editingMilestone ? 'Update milestone' : 'Add milestone'}</button>{editingMilestone && <button type="button" onClick={() => { setEditingMilestone(null); setMilestone(emptyMilestone) }}>Cancel</button>}</form><ul className="admin-list">{project.milestones.map((item) => <li key={item.id}><span><strong>{item.date} · {item.title}</strong><br />{item.description}</span><span className="list-actions"><button type="button" onClick={() => { setEditingMilestone(item.id ?? null); setMilestone({ date: item.date, title: item.title, description: item.description, media_ids: item.media_ids }) }}>Edit</button><button type="button" onClick={() => void deleteHistory(item)}>Delete</button></span></li>)}</ul></section></>}
  </main>
}
