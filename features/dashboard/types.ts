export type Link = { id?: string; project_id?: string; kind: 'github' | 'trello'; url: string; label?: string }
export type Media = { id?: string; project_id?: string; role: 'thumbnail' | 'original' | 'screenshot' | 'video'; source: string; alt_text?: string; caption?: string }
export type Project = { id: string; name: string; description?: string; agentic_platform?: string; technologies: string[]; links: Link[]; media: Media[]; milestones: unknown[] }
