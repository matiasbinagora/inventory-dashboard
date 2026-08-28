export type Link = { id?: string; project_id?: string; kind: 'github' | 'trello'; url: string; label?: string }
export type Media = { id?: string; project_id?: string; role: 'thumbnail' | 'original' | 'screenshot' | 'video'; source: string; original_media_id?: string; alt_text?: string; caption?: string; curated: boolean }
export type Milestone = { id?: string; project_id?: string; date: string; title: string; description: string; media_ids: string[] }
export type Project = { id: string; name: string; description?: string; agentic_platform?: string; github_repository_url?: string; trello_backlog_url?: string; technologies: string[]; links: Link[]; media: Media[]; milestones: Milestone[] }
