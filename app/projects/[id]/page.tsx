import { ProjectDetail } from '../../../features/project-detail/project-detail'

export default async function ProjectPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  return <ProjectDetail projectID={id} />
}
